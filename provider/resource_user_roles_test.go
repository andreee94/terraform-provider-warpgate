package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserRolesResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserRolesResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_user_roles.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserRolesUpdatedResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_user_roles.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserRolesUpdateAddResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_user_roles.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserRolesUpdateRemoveResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_user_roles.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserRolesUpdateRemoveAllResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckFuncValidUUID("warpgate_user_roles.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserRolesResourceConfig() string {
	return `
provider "warpgate" {}
	  
resource "warpgate_user" "one" {
	username = "one"
	credentials = [
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		}
	]
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

resource "warpgate_user_roles" "test" {
	id = warpgate_user.one.id
	role_ids = [
		warpgate_role.one.id,
		warpgate_role.two.id,
	]
}
`
}

func testAccUserRolesUpdatedResourceConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_user" "one" {
	username = "one"
	credentials = [
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		}
	]
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

resource "warpgate_user_roles" "test" {
	id = warpgate_user.one.id
	role_ids = [
		warpgate_role.two.id,
		warpgate_role.one.id,
	]
}
`
}

func testAccUserRolesUpdateAddResourceConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_user" "one" {
	username = "one"
	credentials = [
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		}
	]
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

resource "warpgate_user_roles" "test" {
	id = warpgate_user.one.id
	role_ids = [
		warpgate_role.three.id,
		warpgate_role.one.id,
		warpgate_role.two.id,
	]
}
`
}

func testAccUserRolesUpdateRemoveResourceConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_user" "one" {
	username = "one"
	credentials = [
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		}
	]
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

resource "warpgate_user_roles" "test" {
	id = warpgate_user.one.id
	role_ids = [
		warpgate_role.three.id,
		warpgate_role.one.id,
	]
}
`
}

func testAccUserRolesUpdateRemoveAllResourceConfig() string {
	return `
provider "warpgate" {}

resource "warpgate_user" "one" {
	username = "one"
	credentials = [
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		}
	]
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

resource "warpgate_user_roles" "test" {
	id = warpgate_user.one.id
	role_ids = [
	]
}
`
}
