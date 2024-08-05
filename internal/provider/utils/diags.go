package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	paths "github.com/hashicorp/terraform-plugin-framework/path"
)

type DiagsHandler struct {
	*diag.Diagnostics
	Title string
	Path  paths.Path
}

func NewDiagsHandler(diags *diag.Diagnostics, title string) DiagsHandler {
	return DiagsHandler{
		Diagnostics: diags,
		Title:       title,
		Path:        paths.Empty(),
	}
}

// WithPath
// Returns a new DiagsHandler with the same diagnostics and title, but the
// provided path. All error/warning appending will be attribute-based.
func (d DiagsHandler) WithPath(path paths.Path) DiagsHandler {
	return DiagsHandler{
		Diagnostics: d.Diagnostics,
		Title:       d.Title,
		Path:        path,
	}
}

func (d DiagsHandler) HandleErr(err error) bool {
	if err != nil {
		d.HandleErrMsg(err.Error())
		return true
	}
	return d.HasError()
}

func (d DiagsHandler) HandleErrMsg(msg string, args ...any) {
	detail := fmt.Sprintf(msg, args...)
	if d.Path.Equal(paths.Empty()) {
		d.AddError(d.Title, detail)
	} else {
		d.AddAttributeError(d.Path, d.Title, detail)
	}
}

func (d DiagsHandler) HandleDiags(diags diag.Diagnostics) bool {
	d.Append(diags...)
	return d.HasError()
}

func (d DiagsHandler) HandlePanic(panicData any) {
	if panicData == nil {
		return
	}
	if err, ok := panicData.(error); ok {
		d.HandleErr(fmt.Errorf("code panic: %w", err))
		return
	}
	d.HandleErrMsg("code panic: %s", panicData)
}

// Overrides ///////////////////////////////////////////////////////////////////

// OLD /////////////////////////////////////////////////////////////////////////

func UnwrapDiag[T any](
	target *diag.Diagnostics,
	path paths.Path,
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
	path paths.Path,
	value *string,
) timetypes.RFC3339 {
	return UnwrapDiag(diagnostics, path, func() (timetypes.RFC3339, diag.Diagnostics) {
		return timetypes.NewRFC3339PointerValue(value)
	})
}

func UnwrapObjectValue(
	diagnostics *diag.Diagnostics,
	path paths.Path,
	attributeTypes map[string]attr.Type,
	attributes map[string]attr.Value,
) types.Object {
	return UnwrapDiag(diagnostics, path, func() (types.Object, diag.Diagnostics) {
		return types.ObjectValue(attributeTypes, attributes)
	})
}

func UnwrapMapValue(
	diagnostics *diag.Diagnostics,
	path paths.Path,
	elementType attr.Type,
	elements map[string]attr.Value,
) types.Map {
	return UnwrapDiag(diagnostics, path, func() (types.Map, diag.Diagnostics) {
		return types.MapValue(elementType, elements)
	})
}

func UnwrapSetValue(
	diagnostics *diag.Diagnostics,
	path paths.Path,
	elementType attr.Type,
	elements []attr.Value,
) types.Set {
	return UnwrapDiag(diagnostics, path, func() (types.Set, diag.Diagnostics) {
		return types.SetValue(elementType, elements)
	})
}

func UnwrapListValue(
	diagnostics *diag.Diagnostics,
	path paths.Path,
	elementType attr.Type,
	elements []attr.Value,
) types.List {
	return UnwrapDiag(diagnostics, path, func() (types.List, diag.Diagnostics) {
		return types.ListValue(elementType, elements)
	})
}

func DiagsErrHandler(diags *diag.Diagnostics, title string) func(error) bool {
	return func(err error) bool {
		if err != nil {
			diags.AddError(title, err.Error())
			return true
		}
		return diags.HasError()
	}
}

func DiagsHandleRecover(errHandler func(err error) bool, panicData any) {
	if panicData == nil {
		return
	}
	if err, ok := panicData.(error); ok {
		errHandler(fmt.Errorf("code panic: %w", err))
		return
	}
	errHandler(fmt.Errorf("code panic: %s", panicData))
}

func DiagsValHandler[T any](errHandler func(error)) func(T, error) T {
	return func(val T, err error) T {
		errHandler(err)
		return val
	}
}
