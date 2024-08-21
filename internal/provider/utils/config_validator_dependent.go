package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	paths "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	. "terraform-provider-idmc/internal/utils"
)

var _ resource.ConfigValidator = &DependentValidator{}
var _ datasource.ConfigValidator = &DependentValidator{}

type DependentValidator struct {
	Dependent    paths.Expression
	Dependencies paths.Expressions
}

func (d DependentValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The %s attribute depends on %s", d.Dependent, d.Dependencies)
}

func (d DependentValidator) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d DependentValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, rsp *datasource.ValidateConfigResponse) {
	d.Validate(NewDiagsHandler(ctx, &rsp.Diagnostics, "Failed to validate resource"), req.Config)
}

func (d DependentValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, rsp *resource.ValidateConfigResponse) {
	d.Validate(NewDiagsHandler(ctx, &rsp.Diagnostics, "Failed to validate resource"), req.Config)
}

func (d DependentValidator) Validate(diags DiagsHandler, config tfsdk.Config) {
	diags = diags.WithTitle("Dependent Attribute Validation")

	// Find any matching paths for the dependent.
	dependentPaths := UnwrapDiags(diags, func() (paths.Paths, diag.Diagnostics) {
		return config.PathMatches(diags.Ctx, d.Dependent)
	})

	// Check if any of them are configured.
	var configuredPaths paths.Paths
	for _, dependentPath := range dependentPaths {
		pathDiags := diags.WithPath(dependentPath)

		// Attempt to read the value.
		var value attr.Value
		if pathDiags.Append(config.GetAttribute(pathDiags.Ctx, dependentPath, &value)) {
			return
		}

		// The value is configured.
		if !value.IsNull() && !value.IsUnknown() {
			configuredPaths.Append(dependentPath)
		}

	}

	// If none of the paths are configured, we need to do any more work.
	if len(configuredPaths) == 0 {
		return
	}

	// Collect all the matching paths for the dependencies.
	dependencyPaths := SliceFlatMap(d.Dependencies, func(dependencyExpression paths.Expression) []paths.Path {
		return UnwrapDiags(diags, func() (paths.Paths, diag.Diagnostics) {
			return config.PathMatches(diags.Ctx, dependencyExpression)
		})
	})

	// Check if _any_ of the matched dependencies aren't configured.
	dependencyNotConfigured := false
	for _, dependencyPath := range dependencyPaths {
		pathDiags := diags.WithPath(dependencyPath)

		// Attempt to read the value.
		var value attr.Value
		if pathDiags.Append(config.GetAttribute(pathDiags.Ctx, dependencyPath, &value)) {
			return
		}

		// The value is explicitly not configured.
		if value.IsNull() {
			dependencyNotConfigured = true
			break
		}

	}

	// If all dependencies are configured, we can stop.
	if !dependencyNotConfigured {
		return
	}

	// Raise an error on each of the configured paths.
	for _, configuredPath := range configuredPaths {
		pathDiags := diags.WithPath(configuredPath)
		pathDiags.AddError(d.Description(pathDiags.Ctx))
	}

}
