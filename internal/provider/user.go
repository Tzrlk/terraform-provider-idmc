package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v3 "terraform-provider-idmc/internal/idmc/v3"
	. "terraform-provider-idmc/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
	. "terraform-provider-idmc/internal/provider/utils"
)

type UserResource struct {
	IdmcProviderResource
}

var _ ResourceWithConfigure = &UserResource{}
var _ ResourceWithImportState = &UserResource{}
var _ ResourceWithConfigValidators = &UserResource{}

func NewUserResource() Resource {
	return &UserResource{
		IdmcProviderResource{
			Name: "user",
		},
	}
}

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
}

// Schema <editor-fold desc="Schema" defaultstate="collapsed">
func (r UserResource) Schema(ctx context.Context, req SchemaRequest, rsp *SchemaResponse) {
	r.IdmcProviderResource.Schema(ctx, req, rsp)
	rsp.Schema.MarkdownDescription = `
Use the users resource to request Informatica Intelligent Cloud Services user
details, create users, update role and user group assignments, and delete users.

NOTE: This resource uses a dynamic rate limit. When the system experiences a
large volume or size of requests, responses might be slow or fail with the error
message, "too many requests."

Use the users resource along with the user_group and role resources to manage
user privileges for Informatica Intelligent Cloud Services tasks and assets.
Users and groups can perform tasks and access assets based on the roles that you
assign to them.

Backed by api operations outlined in [the IDMC user api docs][users].

[users]: https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users.html
`
	rsp.Schema.Attributes = MapMerge(rsp.Schema.Attributes, map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "Informatica Intelligent Cloud Services user name. Maximum length is 255 characters.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(255),
			},
		},
		"first_name": schema.StringAttribute{
			Description: "First name for the user account.",
			Required:    true,
		},
		"last_name": schema.StringAttribute{
			Description: "Last name for the user account.",
			Required:    true,
		},
		"password": schema.StringAttribute{
			Description: "Informatica Intelligent Cloud Services password.\nIf password is empty, the user receives an activation email.",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(255),
			},
		},
		"description": schema.StringAttribute{
			Description: "Description of the user.",
			Optional:    true,
		},
		"email": schema.StringAttribute{
			Description: "Email address for the user.",
			Required:    true,
		},
		"title": schema.StringAttribute{
			Description: "Job title of the user.",
			Optional:    true,
		},
		"phone": schema.StringAttribute{
			Description: "Phone number for the user.",
			Optional:    true,
		},
		"force_password_change": schema.BoolAttribute{
			Description: "Whether the user must reset the password after the user logs in for the first time.",
			Optional:    true,
		},
		"max_login_attempts": schema.Int32Attribute{
			Description: "Number of times a user can attempt to log in before the account is locked.",
			Optional:    true,
			Validators: []validator.Int32{
				int32validator.AtLeast(1),
			},
		},
		"auth_mode": schema.StringAttribute{
			Description: "The authentication mode configured by this user (Native/SAML)",
			Optional:    true,
		},
		"auth_alias": schema.StringAttribute{
			Description: "The user identifier or user name in the 3rd party system.",
			Optional:    true,
		},
		"roles": schema.ListAttribute{
			Description: "IDs of the roles to assign to the user.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"groups": schema.ListAttribute{
			Description: "IDs of the user groups to assign to the user.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"last_login_time": schema.StringAttribute{
			Description: "Date and time the user last logged-in.",
			CustomType:  timetypes.RFC3339Type{},
			Computed:    true,
		},
		"last_login_mode": schema.StringAttribute{
			Description: "Whether the user logged in through a REST API call or through the UI.",
			Computed:    true,
		},
		"timezone": schema.StringAttribute{
			MarkdownDescription: `
Time zone of the user.<br/>
For more information, see [Time zone codes](https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/rest-api-codes/time-zone-codes.html).
`,
			Computed: true,
		},
		"state": schema.StringAttribute{
			Computed: true,
			MarkdownDescription: `
State of the user account. Returns one of the following values:<br/>
* Active. User account exists and user has activated the account.
* Provisioned. User account exists but the user has not activated the account.
* Disabled. User account is disabled because the user exceeded the maximum number of login attempts.

NOTE: If the user's password is expired, the value is null.
`,
		},
	})
}

// </editor-fold>

func (r UserResource) ConfigValidators(_ context.Context) []ConfigValidator {
	return []ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("roles"),
			path.MatchRoot("groups"),
		),
		DependentValidator{
			Dependent: path.MatchRoot("auth_alias"),
			Dependencies: path.Expressions{
				path.MatchRoot("auth_mode"),
			},
		},
	}
}

// Create <editor-fold desc="Create" defaultstate="collapsed">
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/creating-a-user.html
func (r UserResource) Create(ctx context.Context, req CreateRequest, rsp *CreateResponse) {
	diags := NewDiagsHandler(ctx, &rsp.Diagnostics, MsgResourceBadCreate)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApi().V3.Client

	var data UserResourceModel
	if diags.Append(req.Plan.Get(ctx, &data)) {
		return
	}

	reqData := v3.CreateUserRequestBody{
		Name:                data.Name.ValueString(),
		FirstName:           data.FirstName.ValueString(),
		LastName:            data.LastName.ValueString(),
		Description:         data.Description.ValueStringPointer(),
		Email:               data.Email.ValueString(),
		Title:               data.Title.ValueStringPointer(),
		Phone:               NullableFromPointer(data.Phone.ValueStringPointer()),
		ForcePasswordChange: data.ForcePasswordChange.ValueBoolPointer(),
		MaxLoginAttempts:    IntPtrFromInt32Attr(data.MaxLoginAttempts),
		AliasName:           data.AuthAlias.ValueStringPointer(),
	}

	// Handle annoying auth edge-case.
	if !data.AuthMode.IsNull() {
		switch data.AuthMode.ValueString() {
		case "Native":
			reqData.Authentication = Ptr(v3.CreateUserRequestBodyAuthenticationN0)
		case "SAML":
			reqData.Authentication = Ptr(v3.CreateUserRequestBodyAuthenticationN1)
		default:
			diags.AddError("Invalid authentication mode: %s", data.AuthMode.ValueString())
		}
	}

	// TODO: roles

	// TODO: groups

	// Actually fire-off the request.
	apiRes, err := client.CreateUserWithResponse(ctx, &v3.CreateUserParams{}, reqData)
	if diags.HandleError(err) {
		return
	}

	// Handle error cases.
	if apiRes.StatusCode != 201 {
		CheckApiErrorV3(diags,
			apiRes.JSON400,
			apiRes.JSON401,
			apiRes.JSON403,
			apiRes.JSON404,
			apiRes.JSON500,
			apiRes.JSON502,
			apiRes.JSON503,
		)
		if !diags.HasError() {
			diags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 201))
		}
		return
	}

	respData := *apiRes.JSON200

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(respData.Id)
	data.Name = types.StringValue(respData.UserName)
	data.Description = types.StringValue(respData.Description)
	data.FirstName = types.StringValue(respData.FirstName)
	data.LastName = types.StringValue(respData.LastName)
	data.Email = types.StringValue(respData.Email)
	data.Phone = NullableToStringAttr(respData.Phone)
	data.Title = types.StringValue(respData.Title)
	data.ForcePasswordChange = types.BoolValue(respData.ForcePasswordChange)
	data.MaxLoginAttempts = types.Int32Value(int32(respData.MaxLoginAttempts))
	data.AuthAlias = types.StringPointerValue(respData.AliasName)

	// Update derived values
	data.OrgId = types.StringPointerValue(respData.OrgId)
	data.State = types.StringPointerValue((*string)(respData.State))
	data.CreatedBy = types.StringPointerValue(respData.CreatedBy)
	data.CreatedTime = diags.TimePointer(respData.CreateTime)
	data.UpdatedBy = types.StringPointerValue(respData.UpdatedBy)
	data.UpdatedTime = diags.TimePointer(respData.UpdateTime)

	// Update annoying values
	switch respData.Authentication {
	case 0:
		data.AuthMode = types.StringValue("Native")
	case 1:
		data.AuthMode = types.StringValue("SAML")
	default:
		data.AuthMode = types.StringUnknown()
	}

	// TODO: roles
	// TODO: groups

	// If we had trouble parsing the data.
	if diags.HasError() {
		return
	}

	// Save creation result back to state.
	diags.Append(rsp.State.Set(ctx, &data))

}

// </editor-fold>

// Read <editor-fold desc="Read" defaultstate="collapsed">
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/getting-user-details.html
func (r UserResource) Read(ctx context.Context, req ReadRequest, rsp *ReadResponse) {
	//TODO implement me
	panic("implement me")
}

// </editor-fold>

// Update <editor-fold desc="Update" defaultstate="collapsed">
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/updating-role-assignments.html
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/updating-user-group-assignments.html
func (r UserResource) Update(ctx context.Context, req UpdateRequest, rsp *UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

// </editor-fold>

// Delete <editor-fold desc="Delete" defaultstate="collapsed">
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/deleting-a-user.html
func (r UserResource) Delete(ctx context.Context, req DeleteRequest, rsp *DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

// </editor-fold>
