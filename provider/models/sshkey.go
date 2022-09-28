package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type SshKey struct {
	Kind            types.String `tfsdk:"kind"`
	PublicKeyBase64 types.String `tfsdk:"public_key_base64"`
}
