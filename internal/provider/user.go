package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

// Schema <editor-fold desc="Schema" defaultstate="collapsed">
func (r UserResource) Schema(ctx context.Context, req SchemaRequest, rsp *SchemaResponse) {
	rsp.Schema = schema.Schema{
		MarkdownDescription: `
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
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
				Description: "Determines whether the user must reset the password after the user logs in for the first time.",
				Optional:    true,
			},
			"max_login_attempts": schema.Int32Attribute{
				Description: "Number of times a user can attempt to log in before the account is locked.",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"saml": schema.BoolAttribute{
				Description: "Determines whether the user accesses Informatica Intelligent Cloud Services through single sign-in (SAML).",
				Optional:    true,
			},
			"saml_alias": schema.StringAttribute{
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
			// vvvvvv from creation response.
			"last_login_time": nil,
			"last_login_mode": nil,
			"timezone":        nil,
			"state":           nil,
			"created_by":      nil,
			"updated_by":      nil,
			"created_time":    nil,
			"updated_time":    nil,
			"org_id":          nil,
		},
	}
}

// </editor-fold>

func (r UserResource) ConfigValidators(ctx context.Context) []ConfigValidator {
	return []ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("roles"),
			path.MatchRoot("groups"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("saml"),
			path.MatchRoot("saml_alias"),
		),
	}
}

// Create <editor-fold desc="Create" defaultstate="collapsed">
// https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-3-resources/users/creating-a-user.html
func (r UserResource) Create(ctx context.Context, req CreateRequest, rsp *CreateResponse) {
	//TODO implement me
	panic("implement me")
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
