package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRoleListDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create roles for testing the datasource
			{
				Config: testAccRolesResourcesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_role.one", "name", "one"),
					testCheckFuncValidUUID("warpgate_role.one", "id"),

					resource.TestCheckResourceAttr("warpgate_role.two", "name", "two"),
					testCheckFuncValidUUID("warpgate_role.two", "id"),

					resource.TestCheckResourceAttr("warpgate_role.three", "name", "three"),
					testCheckFuncValidUUID("warpgate_role.three", "id"),
				),
			},
			// Test the datasource
			{
				Config: testAccRoleListDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.warpgate_role_list.test", "roles.#", "4"),
					testCheckFuncValidUUID("data.warpgate_role_list.test", "roles.0.id"),
					testCheckFuncValidUUID("data.warpgate_role_list.test", "roles.1.id"),
					testCheckFuncValidUUID("data.warpgate_role_list.test", "roles.2.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRolesResourcesConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}
`
}

func testAccRoleListDataSourceConfig() string {
	return `
provider "warpgate" {}

data "warpgate_role_list" "test" {
}
`
}
