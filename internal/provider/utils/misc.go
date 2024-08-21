package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oapi-codegen/nullable"

	. "terraform-provider-idmc/internal/utils"
)

func NullableFromPointer[T any](ptr *T) nullable.Nullable[T] {
	if ptr == nil {
		return nullable.NewNullNullable[T]()
	} else {
		return nullable.NewNullableWithValue[T](*ptr)
	}
}

func NullableToStringAttr(value nullable.Nullable[string]) types.String {
	if value.IsNull() {
		return types.StringNull()
	}
	if !value.IsSpecified() {
		return types.StringUnknown()
	}
	return types.StringValue(value.MustGet())
}

func IntPtrFromInt32Attr(attr types.Int32) *int {
	if attr.IsNull() {
		return nil
	}

	return Ptr(int(attr.ValueInt32()))
}

func TfSetToSlice[T any](diags DiagsHandler, tfSet types.Set) *[]T {
	if tfSet.IsNull() {
		return nil
	}
	var results []T
	diags.Append(tfSet.ElementsAs(diags.Ctx, &results, false))
	return &results
}
