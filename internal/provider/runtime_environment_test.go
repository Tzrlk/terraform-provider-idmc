package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	. "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRuntimeEnvironment(t *testing.T) {
	envName := acctest.RandomWithPrefix("test_env_")
	ParallelTest(t, TestCase{
		ProtoV6ProviderFactories: testAccProviders,
		PreCheck: testAccPreCheck(t),
		Steps: []TestStep{
			{	// Create
				ConfigFile: config.StaticFile("examples/resources/idmc_runtime_environment/resource.tf"),
				ConfigVariables: map[string]config.Variable{
					"name":   config.StringVariable(envName),
					"shared": config.BoolVariable(false),
				},
				Check: ComposeAggregateTestCheckFunc(
					TestCheckResourceAttr("idmc_runtime_environment.example", "name", envName),
					TestCheckResourceAttr("idmc_runtime_environment.example", "shared", "false"),
				),
			},
			{	// Change name
				ConfigFile: config.StaticFile("examples/resources/idmc_runtime_environment/resource.tf"),
				ConfigVariables: map[string]config.Variable{
					"name":   config.StringVariable(envName + "_changed"),
					"shared": config.BoolVariable(false),
				},
				Check: ComposeAggregateTestCheckFunc(
					TestCheckResourceAttr("idmc_runtime_environment.example", "name", envName + "_changed"),
					TestCheckResourceAttr("idmc_runtime_environment.example", "shared", "false"),
				),
			},
			{	// Change shared
				ConfigFile: config.StaticFile("examples/resources/idmc_runtime_environment/resource.tf"),
				ConfigVariables: map[string]config.Variable{
					"name":   config.StringVariable(envName + "_changed"),
					"shared": config.BoolVariable(true),
				},
				Check: ComposeAggregateTestCheckFunc(
					TestCheckResourceAttr("idmc_runtime_environment.example", "name", envName),
					TestCheckResourceAttr("idmc_runtime_environment.example", "shared", "true"),
				),
			},
		},
	})
}
