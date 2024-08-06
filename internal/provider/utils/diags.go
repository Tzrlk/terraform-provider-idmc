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

func (d DiagsHandler) WithTitle(title string) DiagsHandler {
	return DiagsHandler{
		Diagnostics: d.Diagnostics,
		Title:       title,
		Path:        d.Path,
	}
}

func (d DiagsHandler) HandleErrMsg(msg string, args ...any) {
	detail := fmt.Sprintf(msg, args...)
	if d.Path.Equal(paths.Empty()) {
		d.Diagnostics.AddError(d.Title, detail)
	} else {
		d.Diagnostics.AddAttributeError(d.Path, d.Title, detail)
	}
}

func (d DiagsHandler) HandleWarnMsg(msg string, args ...any) {
	detail := fmt.Sprintf(msg, args...)
	if d.Path.Equal(paths.Empty()) {
		d.Diagnostics.AddWarning(d.Title, detail)
	} else {
		d.Diagnostics.AddAttributeWarning(d.Path, d.Title, detail)
	}
}

func (d DiagsHandler) HandleError(err error) bool {
	if err != nil {
		d.HandleErrMsg(err.Error())
		return true
	}
	return d.HasError()
}

func (d DiagsHandler) HandleDiags(diags diag.Diagnostics) bool {

	// If we're at the root, we don't need to cook anything.
	if d.Path.Equal(paths.Empty()) {
		d.Diagnostics.Append(diags...)
		return d.HasError()
	}

	// Otherwise, modify all incoming diags to be pathed.
	for _, diagItem := range diags {
		d.Diagnostics.Append(diag.WithPath(d.Path, diagItem))
	}

	return d.HasError()
}

func (d DiagsHandler) HandlePanic(panicData any) {
	if panicData == nil {
		return
	}
	if err, ok := panicData.(error); ok {
		d.HandleError(fmt.Errorf("code panic: %w", err))
		return
	}
	d.HandleErrMsg("code panic: %s", panicData)
}

// Path Navigation /////////////////////////////////////////////////////////////

func (d DiagsHandler) AtListIndex(index int) DiagsHandler {
	return d.WithPath(d.Path.AtListIndex(index))
}

func (d DiagsHandler) AtMapKey(key string) DiagsHandler {
	return d.WithPath(d.Path.AtMapKey(key))
}

func (d DiagsHandler) AtName(name string) DiagsHandler {
	return d.WithPath(d.Path.AtName(name))
}

func (d DiagsHandler) AtSetValue(value attr.Value) DiagsHandler {
	return d.WithPath(d.Path.AtSetValue(value))
}

func (d DiagsHandler) AtTupleIndex(index int) DiagsHandler {
	return d.WithPath(d.Path.AtTupleIndex(index))
}

// Diags Unwrapping ////////////////////////////////////////////////////////////

func (d DiagsHandler) SetValue(elementType attr.Type, elements []attr.Value) types.Set {
	result, diags := types.SetValue(elementType, elements)
	d.HandleDiags(diags)
	return result
}

func (d DiagsHandler) ListValue(elementType attr.Type, elements []attr.Value) types.List {
	result, diags := types.ListValue(elementType, elements)
	d.HandleDiags(diags)
	return result
}

func (d DiagsHandler) MapValue(elementType attr.Type, elements map[string]attr.Value) types.Map {
	result, diags := types.MapValue(elementType, elements)
	d.HandleDiags(diags)
	return result
}

func (d DiagsHandler) ObjectValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) types.Object {
	result, diags := types.ObjectValue(attributeTypes, attributes)
	d.HandleDiags(diags)
	return result
}

func (d DiagsHandler) TimeValue(text string) timetypes.RFC3339 {
	result, diags := timetypes.NewRFC3339Value(text)
	d.HandleDiags(diags)
	return result
}

func (d DiagsHandler) TimePointer(text *string) timetypes.RFC3339 {
	result, diags := timetypes.NewRFC3339PointerValue(text)
	d.HandleDiags(diags)
	return result
}
