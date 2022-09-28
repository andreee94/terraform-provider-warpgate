package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/bxcodec/faker/v4/pkg/options"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserResource(t *testing.T) {

	type Data struct {
		TotpKey []int8 `faker:"slice_len=32"`
	}

	data := Data{}
	_ = faker.FakeData(&data, options.WithRandomMapAndSliceMaxSize(32)) // If no slice_len is set, this sets the max of the random size

	totp_key, _ := json.Marshal(data.TotpKey)
	totp_key_string := string(totp_key)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig("one", totp_key_string),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_user.test", "username", "one"),
					resource.TestCheckResourceAttr("warpgate_user.test", "credentials.#", "6"),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":     "Sso",
						"email":    "test@example.com",
						"provider": "",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":     "Sso",
						"email":    "test2@example.com",
						"provider": "",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "PublicKey",
						"public_key": "AAAAAAAAAAA",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "PublicKey",
						"public_key": "BBBBBBBBBBB",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind": "Password",
						"hash": "$argon2id$v=19$m=65536,t=1,p=2$5rAIZSCP/YX+JM8m7mo4gQ$TSGk41+4MOzCPbDOjB2AdU18Mz57Df4hmWyNjoilu7k",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "Totp",
						"totp_key.#": "32",
						// "totp_key": fmt.Sprintf("%v", totp_key_string),
					}),

					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "warpgate_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig("two", totp_key_string),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_user.test", "username", "two"),
					resource.TestCheckResourceAttr("warpgate_user.test", "credentials.#", "6"),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":     "Sso",
						"email":    "test@example.com",
						"provider": "",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":     "Sso",
						"email":    "test2@example.com",
						"provider": "",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "PublicKey",
						"public_key": "AAAAAAAAAAA",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "PublicKey",
						"public_key": "BBBBBBBBBBB",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind": "Password",
						"hash": "$argon2id$v=19$m=65536,t=1,p=2$5rAIZSCP/YX+JM8m7mo4gQ$TSGk41+4MOzCPbDOjB2AdU18Mz57Df4hmWyNjoilu7k",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "Totp",
						"totp_key.#": "32",
						// "totp_key": fmt.Sprintf("%v", totp_key_string),
					}),

					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserUpdateRemoveCredentialsResourceConfig("two", totp_key_string),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("warpgate_user.test", "username", "two"),
					resource.TestCheckResourceAttr("warpgate_user.test", "credentials.#", "4"),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":     "Sso",
						"email":    "test@example.com",
						"provider": "",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "PublicKey",
						"public_key": "CCCCCCCCCCC",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind": "Password",
						"hash": "$argon2id$v=19$m=65536,t=1,p=2$5rAIZSCP/YX+JM8m7mo4gQ$TSGk41+4MOzCPbDOjB2AdU18Mz57Df4hmWyNjoilu7k",
					}),

					resource.TestCheckTypeSetElemNestedAttrs("warpgate_user.test", "credentials.*", map[string]string{
						"kind":       "Totp",
						"totp_key.#": "32",
						// "totp_key": fmt.Sprintf("%v", totp_key_string),
					}),

					testCheckFuncValidUUID("warpgate_user.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(name string, totp_key string) string {

	return fmt.Sprintf(`
provider "warpgate" {}
	  
resource "warpgate_user" "test" {
	username = "%s"
	credentials = [
		{
			kind = "Sso"
			email = "test@example.com"
			provider = "" // requires a provider added in the yaml file
		},
		{
			kind = "Sso"
			email = "test2@example.com"
			provider = "" // requires a provider added in the yaml file
		},
		{
			kind = "PublicKey"
			public_key = "AAAAAAAAAAA"
		},
		{
			kind = "PublicKey"
			public_key = "BBBBBBBBBBB"
		},
		{
			kind = "Password"
			hash = "$argon2id$v=19$m=65536,t=1,p=2$5rAIZSCP/YX+JM8m7mo4gQ$TSGk41+4MOzCPbDOjB2AdU18Mz57Df4hmWyNjoilu7k"
		},
		{
			kind = "Totp"
			totp_key = %s
		}
	]
}
`, name, totp_key)
}

func testAccUserUpdateRemoveCredentialsResourceConfig(name string, totp_key string) string {

	return fmt.Sprintf(`
provider "warpgate" {}
	  	  
resource "warpgate_user" "test" {
	username = "%s"
	credentials = [
		{
			kind = "Totp"
			totp_key = %s
		},
		{
			kind = "Sso"
			email = "test@example.com"
			provider = "" // requires a provider added in the yaml file
		},
		{
			kind = "PublicKey"
			public_key = "CCCCCCCCCCC"
		},
		{
			kind = "Password"
			hash = "$argon2id$v=19$m=65536,t=1,p=2$5rAIZSCP/YX+JM8m7mo4gQ$TSGk41+4MOzCPbDOjB2AdU18Mz57Df4hmWyNjoilu7k"
		}
	]
}
`, name, totp_key)
}
