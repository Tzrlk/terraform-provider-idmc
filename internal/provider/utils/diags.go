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
	diags *diag.Diagnostics
	title string
	path  paths.Path
}

func NewDiagsHandler(diags *diag.Diagnostics, title string) DiagsHandler {
	return DiagsHandler{
		diags: diags,
		title: title,
		path:  paths.Empty(),
	}
}

// Overrides ///////////////////////////////////////////////////////////////////

func (d DiagsHandler) HasError() bool {
	return d.diags.HasError()
}

func (d DiagsHandler) AddError(msg string, args ...any) {
	detail := fmt.Sprintf(msg, args...)
	if d.path.Equal(paths.Empty()) {
		d.diags.AddError(d.title, detail)
	} else {
		d.diags.AddAttributeError(d.path, d.title, detail)
	}
}

func (d DiagsHandler) AddWarning(msg string, args ...any) {
	detail := fmt.Sprintf(msg, args...)
	if d.path.Equal(paths.Empty()) {
		d.diags.AddWarning(d.title, detail)
	} else {
		d.diags.AddAttributeWarning(d.path, d.title, detail)
	}
}

func (d DiagsHandler) Append(diags diag.Diagnostics) bool {

	// If we're at the root, we don't need to cook anything.
	if d.path.Equal(paths.Empty()) {
		d.diags.Append(diags...)
		return d.diags.HasError()
	}

	// Otherwise, modify all incoming diags to be pathed.
	for _, diagItem := range diags {
		d.diags.Append(diag.WithPath(d.path, diagItem))
	}

	return d.diags.HasError()
}

// Handling ////////////////////////////////////////////////////////////////////

func (d DiagsHandler) HandleError(err error) bool {
	if err != nil {
		d.AddError(err.Error())
		return true
	}
	return d.diags.HasError()
}

func (d DiagsHandler) HandlePanic(panicData any) {
	if panicData == nil {
		return
	}
	if err, ok := panicData.(error); ok {
		d.HandleError(fmt.Errorf("code panic: %w", err))
		return
	}
	d.AddError("code panic: %s", panicData)
}

// Transformation //////////////////////////////////////////////////////////////

// WithPath
// Returns a new DiagsHandler with the same diagnostics and title, but the
// provided path. All error/warning appending will be attribute-based.
func (d DiagsHandler) WithPath(path paths.Path) DiagsHandler {
	return DiagsHandler{
		diags: d.diags,
		title: d.title,
		path:  path,
	}
}

func (d DiagsHandler) WithTitle(title string) DiagsHandler {
	return DiagsHandler{
		diags: d.diags,
		title: title,
		path:  d.path,
	}
}

// Path Navigation /////////////////////////////////////////////////////////////

func (d DiagsHandler) AtListIndex(index int) DiagsHandler {
	return d.WithPath(d.path.AtListIndex(index))
}

func (d DiagsHandler) AtMapKey(key string) DiagsHandler {
	return d.WithPath(d.path.AtMapKey(key))
}

func (d DiagsHandler) AtName(name string) DiagsHandler {
	return d.WithPath(d.path.AtName(name))
}

func (d DiagsHandler) AtSetValue(value attr.Value) DiagsHandler {
	return d.WithPath(d.path.AtSetValue(value))
}

func (d DiagsHandler) AtTupleIndex(index int) DiagsHandler {
	return d.WithPath(d.path.AtTupleIndex(index))
}

// Diags Unwrapping ////////////////////////////////////////////////////////////

func (d DiagsHandler) SetValue(elementType attr.Type, elements []attr.Value) types.Set {
	result, diags := types.SetValue(elementType, elements)
	d.Append(diags)
	return result
}

func (d DiagsHandler) ListValue(elementType attr.Type, elements []attr.Value) types.List {
	result, diags := types.ListValue(elementType, elements)
	d.Append(diags)
	return result
}

func (d DiagsHandler) MapValue(elementType attr.Type, elements map[string]attr.Value) types.Map {
	result, diags := types.MapValue(elementType, elements)
	d.Append(diags)
	return result
}

func (d DiagsHandler) ObjectValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) types.Object {
	result, diags := types.ObjectValue(attributeTypes, attributes)
	d.Append(diags)
	return result
}

func (d DiagsHandler) TimeValue(text string) timetypes.RFC3339 {
	result, diags := timetypes.NewRFC3339Value(text)
	d.Append(diags)
	return result
}

func (d DiagsHandler) TimePointer(text *string) timetypes.RFC3339 {
	result, diags := timetypes.NewRFC3339PointerValue(text)
	d.Append(diags)
	return result
}
