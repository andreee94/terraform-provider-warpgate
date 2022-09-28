package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHttpTargetResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHttpTargetResourceConfig("one", "10.10.10.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_http_target.test", "name", "one"),
					resource.TestCheckResourceAttr("warpgate_http_target.test", "options.url", "10.10.10.10"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "name"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.url"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.tls.verify"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.tls.mode"),
					resource.TestCheckNoResourceAttr("warpgate_http_target.test", "options.headers"),
					resource.TestCheckNoResourceAttr("warpgate_http_target.test", "options.external_host"),
					testCheckFuncValidUUID("warpgate_http_target.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_http_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccHttpTargetResourceConfig("two", "20.20.20.20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_http_target.test", "name", "two"),
					resource.TestCheckResourceAttr("warpgate_http_target.test", "options.url", "20.20.20.20"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "name"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.url"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.tls.verify"),
					resource.TestCheckResourceAttrSet("warpgate_http_target.test", "options.tls.mode"),
					resource.TestCheckNoResourceAttr("warpgate_http_target.test", "options.headers"),
					resource.TestCheckNoResourceAttr("warpgate_http_target.test", "options.external_host"),
					testCheckFuncValidUUID("warpgate_http_target.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHttpTargetResourceConfig(name string, url string) string {
	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_http_target" "test" {
	name = "%s"
	options = {
		url = "%s"
		tls = {
			mode = "Preferred" 
			verify = true
		}
	}
}
`, name, url)
}
