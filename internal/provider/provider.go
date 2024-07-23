package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-idmc/internal/idmc"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/provider/utils"
)

// Ensure IdmcProvider satisfies various provider interfaces.
var _ provider.Provider = &IdmcProvider{}
var _ provider.ProviderWithFunctions = &IdmcProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IdmcProvider{
			version: version,
		}
	}
}

// IdmcProvider defines the provider implementation.
type IdmcProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// IdmcProviderModel describes the provider data model.
type IdmcProviderModel struct {
	AuthHost types.String `tfsdk:"auth_host"`
	AuthUser types.String `tfsdk:"auth_user"`
	AuthPass types.String `tfsdk:"auth_pass"`
}

type IdmcProviderData struct {
	Api *idmc.IdmcApi
}

func (p *IdmcProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "idmc"
	resp.Version = p.version
}

func (p *IdmcProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "TODO",
		Attributes: map[string]schema.Attribute{
			"auth_host": schema.StringAttribute{
				Description: "The IDMC API authentication host.",
				Optional:    true,
			},
			"auth_user": schema.StringAttribute{
				Description: "The IDMC user name.",
				Optional:    true,
			},
			"auth_pass": schema.StringAttribute{
				Description: "The IDMC user password.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *IdmcProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	diags := &resp.Diagnostics
	errHandler := utils.DiagsErrHandler(diags, "Unable to initialise provider")

	var config IdmcProviderModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	authHost := os.Getenv("IDMC_AUTH_HOST")
	if !config.AuthHost.IsNull() {
		authHost = config.AuthHost.ValueString()
	}
	if authHost == "" {
		diags.AddAttributeError(
			path.Root("auth_host"),
			"Missing IDMC API authentication host",
			"Either 'auth_host' in the config or 'IDMC_AUTH_HOST' in the env is needed.",
		)
	}

	authUser := os.Getenv("IDMC_AUTH_USER")
	if !config.AuthUser.IsNull() {
		authUser = config.AuthUser.ValueString()
	}
	if authUser == "" {
		diags.AddAttributeError(
			path.Root("auth_user"),
			"Missing IDMC API authentication username",
			"Either 'auth_user' in the config or 'IDMC_AUTH_USER' in the env is needed.",
		)
	}

	authPass := os.Getenv("IDMC_AUTH_PASS")
	if !config.AuthPass.IsNull() {
		authPass = config.AuthPass.ValueString()
	}
	if authPass == "" {
		diags.AddAttributeError(
			path.Root("auth_pass"),
			"Missing IDMC API authentication password",
			"Either 'auth_pass' in the config or 'IDMC_AUTH_PASS' in the env is needed.",
		)
	}

	// Catching errors here lets us display all the config issues in one go,
	// rather than treating the user to an onion of failure.
	if diags.HasError() {
		return
	}

	tflog.Debug(ctx, "Setting-up IDMC api client", map[string]any{
		"auth_host": authHost,
		"auth_user": authUser,
	})

	httpClient := &http.Client{}
	baseApiUrl, sessionId, loginErr := doLogin(ctx, authHost, authUser, authPass, httpClient)
	if loginErr != nil {
		errHandler(loginErr)
		return
	}

	idmcApi, idmcApiErr := idmc.NewIdmcApi(baseApiUrl, sessionId,
		common.WithHTTPClient(httpClient),
		common.WithRequestEditorFn(utils.LogHttpRequest),
		common.WithApiResponseEditorFn(utils.LogApiResponse),
	)
	if idmcApiErr != nil {
		errHandler(idmcApiErr)
		return
	}

	idmcProviderData := &IdmcProviderData{
		Api: idmcApi,
	}
	resp.DataSourceData = idmcProviderData
	resp.ResourceData = idmcProviderData

}

func doLogin(ctx context.Context, authHost string, authUser string, authPass string, httpClient common.HttpRequestDoer) (string, string, error) {
	var apiUrl = fmt.Sprintf("https://%s/saas", authHost)

	// First set up a client configured for api login (without logging requests).
	client, clientErr := v3.NewClientWithResponses(apiUrl,
		common.WithHTTPClient(httpClient),
		common.WithRequestEditorFn(func(httpCtx context.Context, req *http.Request) error {
			req.Header["Accept"] = []string{"application/json"}
			return nil
		}),
		common.WithApiResponseEditorFn(utils.LogApiResponse),
	)
	if clientErr != nil {
		return apiUrl, "", clientErr
	}

	// Perform the login operation with the provided credentials.
	res, resErr := client.LoginWithResponse(ctx, v3.LoginJSONRequestBody{
		Username: authUser,
		Password: authPass,
	})
	if resErr != nil {
		return apiUrl, "", resErr
	}

	// We only want 200 responses.
	if err := utils.RequireHttpStatus(200, res); err != nil {
		return apiUrl, "", err
	}

	// Extract the key information from the login response
	if res.JSON200 == nil {
		return apiUrl, "", fmt.Errorf("response data has not been parsed")
	}
	resData := *res.JSON200

	if resData.UserInfo == nil {
		return apiUrl, "", fmt.Errorf("no user data found in response")
	}
	userData := *resData.UserInfo

	if userData.SessionId == nil {
		return apiUrl, "", fmt.Errorf("no sessionId found in response")
	}
	sessionId := *userData.SessionId

	if resData.Products == nil {
		return apiUrl, sessionId, fmt.Errorf("no products found in response")
	}
	products := *resData.Products

	for _, product := range products {
		if product.Name != nil && *product.Name == "Integration Cloud" {
			apiUrl = *product.BaseApiUrl
			return apiUrl, sessionId, nil
		}
	}

	return apiUrl, sessionId, fmt.Errorf("no api url found in response")

}

func (p *IdmcProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRoleResource,
	}
}

func (p *IdmcProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgentInstallerDataSource,
		NewRolesDataSource,
	}
}

func (p *IdmcProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
