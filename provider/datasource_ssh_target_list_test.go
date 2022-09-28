package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSshTargetListDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create roles for testing the datasource
			{
				Config: testAccSshTargetResourcesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.one", "name", "one"),
					testCheckFuncValidUUID("warpgate_ssh_target.one", "id"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.one", "options.host", "10.10.10.10"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.one", "options.port", "11"),

					resource.TestCheckResourceAttr("warpgate_ssh_target.two", "name", "two"),
					testCheckFuncValidUUID("warpgate_ssh_target.two", "id"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.two", "options.host", "20.20.20.20"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.two", "options.port", "22"),

					resource.TestCheckResourceAttr("warpgate_ssh_target.three", "name", "three"),
					testCheckFuncValidUUID("warpgate_ssh_target.three", "id"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.three", "options.host", "30.30.30.30"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.three", "options.port", "33"),
				),
			},
			// Test the datasource
			{
				Config: testAccSshTargetListDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.warpgate_ssh_target_list.test", "targets.#", "3"),

					resource.TestCheckResourceAttr("data.warpgate_ssh_target_list.test", "targets.0.options.auth_kind", "Password"),
					resource.TestCheckResourceAttr("data.warpgate_ssh_target_list.test", "targets.0.options.password", "A12345678"),

					resource.TestCheckResourceAttr("data.warpgate_ssh_target_list.test", "targets.1.options.auth_kind", "PublicKey"),
					resource.TestCheckNoResourceAttr("data.warpgate_ssh_target_list.test", "targets.1.options.password"),

					testCheckFuncValidUUID("data.warpgate_ssh_target_list.test", "targets.0.id"),
					testCheckFuncValidUUID("data.warpgate_ssh_target_list.test", "targets.1.id"),
					testCheckFuncValidUUID("data.warpgate_ssh_target_list.test", "targets.2.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSshTargetResourcesConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 11
		username = "root"
		auth_kind = "Password"
		password = "A12345678"
	}
}
resource "warpgate_ssh_target" "two" {
	name = "two"
	options = {
		host = "20.20.20.20"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}
resource "warpgate_ssh_target" "three" {
	name = "three"
	options = {
		host = "30.30.30.30"
		port = 33
		username = "root"
		auth_kind = "PublicKey"
	}
}
`
}

func testAccSshTargetListDataSourceConfig() string {
	return `
provider "warpgate" {}

data "warpgate_ssh_target_list" "test" {
}
`
}
