package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserResource(t *testing.T) {

	// os.Setenv("TF_ACC", "1")
	// os.Setenv("TF_LOG", "debug")

	// os.Setenv("WARPGATE_HOST", "127.0.0.1")
	// os.Setenv("WARPGATE_PORT", "38888")
	// os.Setenv("WARPGATE_USERNAME", "admin")
	// os.Setenv("WARPGATE_PASSWORD", "password")
	// os.Setenv("WARPGATE_INSECURE_SKIP_VERIFY", "true")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_user.test", "username", "one"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.host", "10.10.10.10"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "id"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "name"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.host"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.port"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.username"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.auth_kind"),
					// resource.TestCheckNoResourceAttr("warpgate_user.test", "options.password"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.auth_kind", "PublicKey"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.username", "root"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.port", "22"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_user.test", "username", "two"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.host", "20.20.20.20"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "id"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.host"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.port"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.username"),
					// resource.TestCheckResourceAttrSet("warpgate_user.test", "options.auth_kind"),
					// resource.TestCheckNoResourceAttr("warpgate_user.test", "options.password"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.auth_kind", "PublicKey"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.username", "root"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "options.port", "22"),
					// resource.TestCheckResourceAttr("warpgate_user.test", "allow_roles.#", "0"),
					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_user" "test" {
	username = "%s"
	credentials = [
		{
			kind = "Sso"
			email = "test@example.com"
			provider = "oidc"
		},
		{
			kind = "Password"
			hash = "password"
		},
		// {
		// 	kind = "Totp"
		// 	topt_key = [0, 1, 2, 3]
		// }
		// ,{
		// 	kind = "PublicKey"
		// 	public_key = "AAAAAAA"
		// }
	]
}
`, name)
}
