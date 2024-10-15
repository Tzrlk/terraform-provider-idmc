package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-idmc/internal/idmc"
	"terraform-provider-idmc/internal/idmc/common"
	. "terraform-provider-idmc/internal/provider/utils"
)

const MsgProviderBadConfigure = "Unable to configure provider"

// Ensure IdmcProvider satisfies various provider interfaces.
var _ provider.Provider = &IdmcProvider{}
var _ provider.ProviderWithFunctions = &IdmcProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IdmcProvider{
			version: version,
			IdmcProviderData: &IdmcProviderData{
				Api: nil,
			},
		}
	}
}

// IdmcProvider defines the provider implementation.
type IdmcProvider struct {
	*IdmcProviderData
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

func (p *IdmcProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "idmc"
	resp.Version = p.version
}

func (p *IdmcProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform_rest_api_version_3_resources/login_2.html",
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

func getCfgVal(diags DiagsHandler, attrVal types.String, attrPath string, required bool) string {

	// Check the attribute for a valid value.
	if !attrVal.IsNull() && attrVal.ValueString() != "" {
		return attrVal.ValueString()
	}

	// Check the environment for a valid value.
	envKey := "IDMC_" + strings.ToUpper(attrPath)
	val, ok := os.LookupEnv(envKey)
	if ok && val != "" {
		return val
	}

	// Register an error on the attribute if required.
	if required {
		diags.AtName(attrPath).AddError(
			"Either '%s' in the config, or '%s' in the env is needed.",
			attrPath, envKey,
		)
	}

	return ""

}

func (p *IdmcProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgProviderBadConfigure)

	// Load the provider config
	var config IdmcProviderModel
	diags.Append(req.Config.Get(ctx, &config))
	if diags.HasError() {
		return
	}

	// Don't re-create the api if it already exists.
	if p.Api != nil {
		tflog.Debug(ctx, "Re-using previously configured api.")
		resp.DataSourceData = p.IdmcProviderData
		resp.ResourceData = p.IdmcProviderData
		return
	}

	// Extract config and validate all the required values are set.
	authHost := getCfgVal(diags, config.AuthHost, "auth_host", true)
	authUser := getCfgVal(diags, config.AuthUser, "auth_user", true)
	authPass := getCfgVal(diags, config.AuthPass, "auth_pass", true)
	if diags.HasError() {
		return
	}

	// Initialise a new IDMC api client.
	api := idmc.NewIdmcApi(
		fmt.Sprintf("https://%s/saas/", authHost),
		common.ClientOptions{
			common.WithRequestEditorFn(LogHttpRequest),
			common.WithApiResponseEditorFn(LogApiResponse),
		},
	)
	if diags.HandleError(api.Error()) {
		return
	}

	// Perform a login to get the api to be functional.
	err := api.MustGet().V3.Login(ctx, authUser, authPass)
	if diags.HandleError(err) {
		return
	}

	// Save the api and set the provider data
	p.Api = api.MustGet()
	resp.DataSourceData = p.IdmcProviderData
	resp.ResourceData = p.IdmcProviderData

}

func (p *IdmcProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRoleResource,
		NewRuntimeEnvironmentResource,
	}
}

func (p *IdmcProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgentInstallerDataSource,
		NewRoleDataSource,
		NewRoleListDataSource,
		NewRolePrivilegeListDataSource,
	}
}

func (p *IdmcProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
