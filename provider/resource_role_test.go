package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRoleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_role.test", "name", "one"),
					resource.TestCheckResourceAttrSet("warpgate_role.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_role.test", "name"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRoleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_role.test", "name", "two"),
					resource.TestCheckResourceAttrSet("warpgate_role.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_role.test", "name"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRoleResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_role" "test" {
	name = "%s"
}
`, name)
}
