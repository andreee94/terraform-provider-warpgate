package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTargetRolesResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTargetRolesResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_target_roles.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_target_roles.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// // Update and Read testing
			{
				Config: testAccTargetRolesUpdatedResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_target_roles.test", "id"),
				),
			},
			// // Update and Read testing
			{
				Config: testAccTargetRolesUpdateAddResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_target_roles.test", "id"),
				),
			},
			// // Update and Read testing
			{
				Config: testAccTargetRolesUpdateRemoveResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_target_roles.test", "id"),
				),
			},
			// // Update and Read testing
			{
				Config: testAccTargetRolesUpdateRemoveAllResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_target_roles.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTargetRolesResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}

resource "warpgate_target_roles" "test" {
	id = warpgate_ssh_target.one.id
	role_ids = [
		warpgate_role.one.id,
		warpgate_role.two.id,
	]
}
`
}

func testAccTargetRolesUpdatedResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}

resource "warpgate_target_roles" "test" {
	id = warpgate_ssh_target.one.id
	role_ids = [
		warpgate_role.two.id,
		warpgate_role.one.id,
	]
}
`
}

func testAccTargetRolesUpdateAddResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}

resource "warpgate_target_roles" "test" {
	id = warpgate_ssh_target.one.id
	role_ids = [
		warpgate_role.three.id,
		warpgate_role.one.id,
		warpgate_role.two.id,
	]
}
`
}

func testAccTargetRolesUpdateRemoveResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}

resource "warpgate_target_roles" "test" {
	id = warpgate_ssh_target.one.id
	role_ids = [
		warpgate_role.three.id,
		warpgate_role.one.id,
	]
}
`
}

func testAccTargetRolesUpdateRemoveAllResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_ssh_target" "one" {
	name = "one"
	options = {
		host = "10.10.10.10"
		port = 22
		username = "root"
		auth_kind = "PublicKey"
	}
}

resource "warpgate_role" "one" {
	name = "one"
}

resource "warpgate_role" "two" {
	name = "two"
}

resource "warpgate_role" "three" {
	name = "three"
}

resource "warpgate_target_roles" "test" {
	id = warpgate_ssh_target.one.id
	role_ids = [
	]
}
`
}
