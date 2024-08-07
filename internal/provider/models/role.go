package models

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RoleType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                  types.StringType,
		"name":                types.StringType,
		"description":         types.StringType,
		"privileges":          RolePrivilegeListType,
		"display_name":        types.StringType,
		"display_description": types.StringType,
		"org_id":              types.StringType,
		"system_role":         types.BoolType,
		"status":              types.StringType,
		"created_by":          types.StringType,
		"updated_by":          types.StringType,
		"created_time":        timetypes.RFC3339Type{},
		"updated_time":        timetypes.RFC3339Type{},
	},
}

type RoleValue struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Privileges         types.List   `tfsdk:"privileges"`
	OrgId              types.String `tfsdk:"org_id"`
	DisplayName        types.String `tfsdk:"display_name"`
	DisplayDescription types.String `tfsdk:"display_description"`
	SystemRole         types.Bool   `tfsdk:"system_role"`
	Status             types.String `tfsdk:"status"`
	CreatedBy          types.String `tfsdk:"created_by"`
	UpdatedBy          types.String `tfsdk:"updated_by"`
	CreatedTime        types.String `tfsdk:"created_time"`
	UpdatedTime        types.String `tfsdk:"updated_time"`
}
