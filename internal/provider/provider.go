package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	AuthUser types.String `tfsdk:"username"`
	AuthPass types.String `tfsdk:"password"`
}

type IdmcProviderData struct {
	Client    *http.Client
	BaseUrl   types.String
	SessionId types.String
}

func (p *IdmcProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "idmc"
	resp.Version = p.version
}

func (p *IdmcProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auth_host": schema.StringAttribute{
				MarkdownDescription: "The IDMC API authentication host.",
				Optional:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "The IDMC user name.",
				Optional:            true,
			},
			"auth_pass": schema.StringAttribute{
				MarkdownDescription: "The IDMC user password.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *IdmcProvider) Configure(
		ctx context.Context,
		req provider.ConfigureRequest,
		resp *provider.ConfigureResponse) {
	var config IdmcProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authHost := os.Getenv("IDMC_AUTH_HOST")
	if !config.AuthHost.IsNull() {
		authHost = config.AuthHost.ValueString()
	}
	if authHost == "" {
		resp.Diagnostics.AddAttributeError(
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
		resp.Diagnostics.AddAttributeError(
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
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_pass"),
			"Missing IDMC API authentication password",
			"Either 'auth_pass' in the config or 'IDMC_AUTH_PASS' in the env is needed.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

}

func (p *IdmcProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
	}
}

func (p *IdmcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	}
}

func (p *IdmcProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
	}
}
