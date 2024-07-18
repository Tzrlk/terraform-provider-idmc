package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/config"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRole(t *testing.T) {
	roleName := acctest.RandString(8)
	roleDesc := acctest.RandString(8)
	privId := acctest.RandString(8) // FIXME: This won't hold up.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("idmc_role.example", "role_name", roleName),
					resource.TestCheckResourceAttr("idmc_role.example", "role_description", roleDesc),
				),
				ConfigDirectory: config.StaticDirectory("../../examples/resources/idmc_role"),
				ConfigVariables: map[string]config.Variable{
					"role_name":         config.StringVariable(roleName),
					"role_description":  config.StringVariable(roleDesc),
					"role_privilege_id": config.StringVariable(privId),
				},
			},
		},
	})
}
