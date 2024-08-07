package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/provider/utils"

	. "terraform-provider-idmc/internal/utils"
)

var RolePrivilegeListType = types.ListType{
	ElemType: RolePrivilegeType,
}

type RolePrivilegeListValue []RolePrivilegeValue

func (r RolePrivilegeListValue) FromApi(diags utils.DiagsHandler, items *[]v3.RolePrivilegeItem) RolePrivilegeListValue {

	if items == nil {
		diags.WithTitle("Issue handling resource API response").AddWarning(
			"Expected role privilege data, but received nothing.")
		return nil
	}

	return TransformSlice(*items, func(from v3.RolePrivilegeItem) RolePrivilegeValue {
		return RolePrivilegeValue{
			Id:          types.StringPointerValue(from.Id),
			Name:        types.StringPointerValue(from.Name),
			Description: types.StringPointerValue(from.Description),
			Service:     types.StringPointerValue(from.Service),
			Status:      types.StringPointerValue((*string)(from.Status)),
		}
	})
}

func NewRolePrivilegeListValueFromList(diags utils.DiagsHandler, list types.List) RolePrivilegeListValue {

	var result RolePrivilegeListValue
	if diags.Append(list.ElementsAs(diags.Ctx, &result, false)) {
		return nil
	}

	return result
}

func (r RolePrivilegeListValue) ToObjectList(diags utils.DiagsHandler) types.List {
	return diags.ListValueFrom(RolePrivilegeType, r)
}

func (r RolePrivilegeListValue) GetIds(diags utils.DiagsHandler) *HashSet[string] {
	return r.extractStrings(diags, func(item RolePrivilegeValue) types.String {
		return item.Id
	})
}

func (r RolePrivilegeListValue) GetNames(diags utils.DiagsHandler) *HashSet[string] {
	return r.extractStrings(diags, func(item RolePrivilegeValue) types.String {
		return item.Name
	})
}

func (r RolePrivilegeListValue) extractStrings(diags utils.DiagsHandler, extractor func(item RolePrivilegeValue) types.String) *HashSet[string] {
	return NewHashSetAfter[string](func(set *HashSet[string]) {
		for index, item := range r {
			itemDiags := diags.AtListIndex(index)
			itemValue := extractor(item)
			if itemValue.IsUnknown() {
				itemDiags.AddWarning("Item not expected to be unknown.")
				continue
			}
			if itemValue.IsNull() {
				itemDiags.AddWarning("Item not expected to be null.")
				continue
			}
			set.Add(itemValue.ValueString())
		}
	})
}
