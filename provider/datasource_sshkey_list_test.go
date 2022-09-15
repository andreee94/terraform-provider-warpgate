package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSshKeyListDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test the datasource
			{
				Config: testAccSshKeyListED25519DataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.ed25519", "sshkeys.#", "1"),
					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.ed25519", "kind", "ssh-ed25519"),
					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.ed25519", "sshkeys.0.kind", "ssh-ed25519"),

					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.rsa_sha2_256", "sshkeys.#", "1"),
					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.rsa_sha2_256", "kind", "rsa-sha2-256"),
					resource.TestCheckResourceAttr("data.warpgate_sshkey_list.rsa_sha2_256", "sshkeys.0.kind", "rsa-sha2-256"),
				),
			},
		},
	})
}

func testAccSshKeyListED25519DataSourceConfig() string {
	return `
provider "warpgate" {}

data "warpgate_sshkey_list" "ed25519" {
	kind = "ssh-ed25519"
}

data "warpgate_sshkey_list" "rsa_sha2_256" {
	kind = "rsa-sha2-256"
}
`
}
