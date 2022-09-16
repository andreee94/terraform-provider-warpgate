package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHttpTargetListDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create roles for testing the datasource
			{
				Config: testAccHttpTargetResourcesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_http_target.one", "name", "one"),
					testCheckFuncValidUUID("warpgate_http_target.one", "id"),
					resource.TestCheckResourceAttr("warpgate_http_target.one", "options.url", "10.10.10.10"),
					// resource.TestCheckResourceAttr("warpgate_http_target.one", "options.port", "11"),

					resource.TestCheckResourceAttr("warpgate_http_target.two", "name", "two"),
					testCheckFuncValidUUID("warpgate_http_target.two", "id"),
					resource.TestCheckResourceAttr("warpgate_http_target.two", "options.url", "20.20.20.20"),
					// resource.TestCheckResourceAttr("warpgate_http_target.two", "options.port", "22"),

					resource.TestCheckResourceAttr("warpgate_http_target.three", "name", "three"),
					testCheckFuncValidUUID("warpgate_http_target.three", "id"),
					resource.TestCheckResourceAttr("warpgate_http_target.three", "options.url", "30.30.30.30"),
					// resource.TestCheckResourceAttr("warpgate_http_target.three", "options.port", "33"),
				),
			},
			// Test the datasource
			{
				Config: testAccHttpTargetListDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 4 instead of 3 since it includes the default warpgate admin page
					resource.TestCheckResourceAttr("data.warpgate_http_target_list.test", "targets.#", "4"),

					resource.TestCheckResourceAttr("data.warpgate_http_target_list.test", "targets.0.options.tls.mode", "Preferred"),
					resource.TestCheckResourceAttr("data.warpgate_http_target_list.test", "targets.1.options.tls.mode", "Preferred"),
					resource.TestCheckResourceAttr("data.warpgate_http_target_list.test", "targets.2.options.tls.mode", "Preferred"),

					// resource.TestCheckResourceAttr("data.warpgate_http_target_list.test", "targets.1.options.auth_kind", "PublicKey"),

					testCheckFuncValidUUID("data.warpgate_http_target_list.test", "targets.0.id"),
					testCheckFuncValidUUID("data.warpgate_http_target_list.test", "targets.1.id"),
					testCheckFuncValidUUID("data.warpgate_http_target_list.test", "targets.2.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHttpTargetResourcesConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_http_target" "one" {
	name = "one"
	options = {
		url = "10.10.10.10"
		tls = {
			mode = "Preferred" 
			verify = true
		}
	}
}
resource "warpgate_http_target" "two" {
	name = "two"
	options = {
		url = "20.20.20.20"
		tls = {
			mode = "Preferred" 
			verify = true
		}
	}
}
resource "warpgate_http_target" "three" {
	name = "three"
	options = {
		url = "30.30.30.30"
		tls = {
			mode = "Preferred" 
			verify = true
		}
	}
}
`
}

func testAccHttpTargetListDataSourceConfig() string {
	return `
provider "warpgate" {}

data "warpgate_http_target_list" "test" {
}
`
}
