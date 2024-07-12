package utils

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
) types.Object {
	return UnwrapDiag(diagnostics, path, func() (types.Object, diag.Diagnostics) {
		return types.ObjectValue(attributeTypes, attributes)
	})
}

func UnwrapMapValue(
	diagnostics *diag.Diagnostics,
	path path.Path,
	elementType attr.Type,
	elements map[string]attr.Value,
) types.Map {
	return UnwrapDiag(diagnostics, path, func() (types.Map, diag.Diagnostics) {
		return types.MapValue(elementType, elements)
	})
}

func UnwrapSetValue(
	diagnostics *diag.Diagnostics,
	path path.Path,
	elementType attr.Type,
	elements []attr.Value,
) types.Set {
	return UnwrapDiag(diagnostics, path, func() (types.Set, diag.Diagnostics) {
		return types.SetValue(elementType, elements)
	})
}
