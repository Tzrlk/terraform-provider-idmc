package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func UnwrapDiag[T any](
	target *diag.Diagnostics,
	path path.Path,
	operation func() (T, diag.Diagnostics),
) T {
	result, source := operation()

	// Append any and all issues to the main diagnostics array.
	for _, diagItem := range source {
		target.Append(diag.WithPath(path, diagItem))
	}

	return result

}

func UnwrapNewRFC3339PointerValue(
	diagnostics *diag.Diagnostics,
	path path.Path,
	value *string,
) timetypes.RFC3339 {
	return UnwrapDiag(diagnostics, path, func() (timetypes.RFC3339, diag.Diagnostics) {
		return timetypes.NewRFC3339PointerValue(value)
	})
}

func UnwrapObjectValue(
	diagnostics *diag.Diagnostics,
	path path.Path,
	attributeTypes map[string]attr.Type,
	attributes map[string]attr.Value,
) basetypes.ObjectValue {
	return UnwrapDiag(diagnostics, path, func() (basetypes.ObjectValue, diag.Diagnostics) {
		return types.ObjectValue(attributeTypes, attributes)
	})
}

func UnwrapMapValue(
	diagnostics *diag.Diagnostics,
	path path.Path,
	elementType attr.Type,
	elements map[string]attr.Value,
) basetypes.MapValue {
	return UnwrapDiag(diagnostics, path, func() (basetypes.MapValue, diag.Diagnostics) {
		return types.MapValue(elementType, elements)
	})
}
