package models

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"

	. "terraform-provider-idmc/internal/provider/utils"
	. "terraform-provider-idmc/internal/utils"
)

type UserResourceModel struct {
	Id                  types.String      `tfsdk:"id"`
	OrgId               types.String      `tfsdk:"org_id"`
	Name                types.String      `tfsdk:"name"`
	Description         types.String      `tfsdk:"description"`
	FirstName           types.String      `tfsdk:"first_name"`
	LastName            types.String      `tfsdk:"last_name"`
	Title               types.String      `tfsdk:"title"`
	Phone               types.String      `tfsdk:"phone"`
	Email               types.String      `tfsdk:"email"`
	State               types.String      `tfsdk:"state"`
	TimeZone            types.String      `tfsdk:"time_zone"`
	MaxLoginAttempts    types.Int32       `tfsdk:"max_login_attempts"`
	AuthMode            types.String      `tfsdk:"auth_mode"`
	AuthAlias           types.String      `tfsdk:"auth_alias"`
	ForcePasswordChange types.Bool        `tfsdk:"force_password_change"`
	LastLoginTime       timetypes.RFC3339 `tfsdk:"last_login_time"`
	LastLoginMode       types.String      `tfsdk:"last_login_mode"`
	CreatedBy           types.String      `tfsdk:"created_by"`
	UpdatedBy           types.String      `tfsdk:"updated_by"`
	CreatedTime         timetypes.RFC3339 `tfsdk:"created_time"`
	UpdatedTime         timetypes.RFC3339 `tfsdk:"updated_time"`
	Roles               types.Set         `tfsdk:"roles"`
	Groups              types.Set         `tfsdk:"groups"`
}

func (u *UserResourceModel) Update(diags DiagsHandler, data *v3.UserDetails) bool {
	if data == nil {
		diags.AddError("Missing parsed API response payload.")
		return true
	}

	// Update the configured state so instabilities can be detected.
	u.Name = types.StringValue(data.UserName)
	u.Description = types.StringValue(data.Description)
	u.FirstName = types.StringValue(data.FirstName)
	u.LastName = types.StringValue(data.LastName)
	u.Email = types.StringValue(data.Email)
	u.Phone = NullableToStringAttr(data.Phone)
	u.Title = types.StringValue(data.Title)
	u.ForcePasswordChange = types.BoolValue(data.ForcePasswordChange)
	u.MaxLoginAttempts = types.Int32Value(int32(data.MaxLoginAttempts))
	u.AuthAlias = types.StringPointerValue(data.AliasName)

	// Update derived values
	u.Id = types.StringPointerValue(data.Id)
	u.OrgId = types.StringPointerValue(data.OrgId)
	u.State = types.StringPointerValue((*string)(data.State))
	u.CreatedBy = types.StringPointerValue(data.CreatedBy)
	u.CreatedTime = diags.TimePointer(data.CreateTime)
	u.UpdatedBy = types.StringPointerValue(data.UpdatedBy)
	u.UpdatedTime = diags.TimePointer(data.UpdateTime)

	// Handle more annoying values
	u.SetAuthMode(data.Authentication)
	u.SetRoles(diags, &data.Roles)
	u.SetGroups(diags, &data.Groups)

	return diags.HasError()

}

func (u *UserResourceModel) GetAuthMode() *v3.UserDetailsAuthentication {
	if u.AuthMode.IsNull() {
		return nil
	}
	switch u.AuthMode.ValueString() {
	case "Native":
		return Ptr(v3.UserDetailsAuthenticationN0)
	case "SAML":
		return Ptr(v3.UserDetailsAuthenticationN1)
	default:
		return nil
	}
}

func (u *UserResourceModel) GetAuthModeForCreate() *v3.CreateUserRequestBodyAuthentication {
	return (*v3.CreateUserRequestBodyAuthentication)(u.GetAuthMode())
}

func (u *UserResourceModel) SetAuthMode(value v3.UserDetailsAuthentication) {
	switch value {
	case v3.UserDetailsAuthenticationN0:
		u.AuthMode = types.StringValue("Native")
	case v3.UserDetailsAuthenticationN1:
		u.AuthMode = types.StringValue("SAML")
	default:
		u.AuthMode = types.StringUnknown()
	}
}

func (u *UserResourceModel) GetRoles(diags DiagsHandler) *[]string {
	return TfSetToSlice[string](diags, u.Roles)
}

func (u *UserResourceModel) SetRoles(diags DiagsHandler, values *[]v3.UserDetailsRole) {
	if values == nil {
		u.Roles = types.SetNull(types.StringType)
		return
	}
	u.Roles = diags.SetValue(types.StringType, TransformSlice(*values, func(from v3.UserDetailsRole) attr.Value {
		return types.StringValue(from.Id)
	}))
}

func (u *UserResourceModel) GetGroups(diags DiagsHandler) *[]string {
	return TfSetToSlice[string](diags, u.Roles)
}

func (u *UserResourceModel) SetGroups(diags DiagsHandler, values *[]v3.UserDetailsGroup) {
	if values == nil {
		u.Roles = types.SetNull(types.StringType)
		return
	}
	u.Groups = diags.SetValue(types.StringType, TransformSlice(*values, func(from v3.UserDetailsGroup) attr.Value {
		return types.StringValue(from.Id)
	}))
}
