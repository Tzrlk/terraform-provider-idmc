package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var RolePrivilegeType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"service":     types.StringType,
		"status":      types.StringType,
	},
}

////////////////////////////////////////////////////////////////////////////////
//var _ basetypes.ObjectValuable = &RolePrivilegeValue{}

type RolePrivilegeValue struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Service     types.String `tfsdk:"service"`
	Status      types.String `tfsdk:"status"`
}

func (r RolePrivilegeValue) Type(_ context.Context) attr.Type {
	return RolePrivilegeType
}

//func (r RolePrivilegeValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
//	return tftypes.New, nil
//}

//func (r RolePrivilegeValue) Equal(value attr.Value) bool {
//	other, ok := value.(types.Object)
//}

func (r RolePrivilegeValue) IsNull() bool {
	return false
}

func (r RolePrivilegeValue) IsUnknown() bool {
	return false
}

//func (r RolePrivilegeValue) String() string {
//}

func (r RolePrivilegeValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, RolePrivilegeType.AttrTypes, r)
}
