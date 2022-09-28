package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSshTargetPublicKeyResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSshTargetPublicKeyResourceConfig("one", "10.10.10.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "one"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "10.10.10.10"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "name"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckNoResourceAttr("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "PublicKey"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_ssh_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSshTargetPublicKeyResourceConfig("two", "20.20.20.20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "two"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "20.20.20.20"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckNoResourceAttr("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "PublicKey"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSshTargetPasswordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSshTargetPasswordResourceConfig("one", "10.10.10.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "one"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "10.10.10.10"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "name"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "Password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.password", "A12345678"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_ssh_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSshTargetPasswordResourceConfig("two", "20.20.20.20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "two"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "20.20.20.20"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "Password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.password", "A12345678"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSshTargetMixedAuthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSshTargetPasswordResourceConfig("one", "10.10.10.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "one"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "10.10.10.10"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "name"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "Password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.password", "A12345678"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_ssh_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSshTargetPublicKeyResourceConfig("two", "20.20.20.20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "name", "two"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.host", "20.20.20.20"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "id"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.host"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.port"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.username"),
					resource.TestCheckResourceAttrSet("warpgate_ssh_target.test", "options.auth_kind"),
					resource.TestCheckNoResourceAttr("warpgate_ssh_target.test", "options.password"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.auth_kind", "PublicKey"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.username", "root"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "options.port", "22"),
					resource.TestCheckResourceAttr("warpgate_ssh_target.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_ssh_target.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSshTargetPublicKeyResourceConfig(name string, host string) string {
	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "test" {
	name = "%s"
	options = {
		host = "%s"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}
`, name, host)
}

func testAccSshTargetPasswordResourceConfig(name string, host string) string {
	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "test" {
	name = "%s"
	options = {
		host = "%s"
		port = 22
		username = "root"
		auth_kind = "Password"
		password = "A12345678"
	}
}
`, name, host)
}
