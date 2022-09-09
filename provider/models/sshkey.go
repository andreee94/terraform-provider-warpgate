package models

type SshKey struct {
	Kind            string `tfsdk:"kind"`
	PublicKeyBase64 string `tfsdk:"public_key_base64"`
}
