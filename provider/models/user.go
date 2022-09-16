package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type User struct {
	// CredentialPolicy *UserRequireCredentialsPolicy `json:"credential_policy,omitempty"`
	Id          types.String         `tfsdk:"id"`
	Username    string               `tfsdk:"username"`
	Credentials []UserAuthCredential `tfsdk:"credentials"`
	Roles       types.Set            `tfsdk:"roles"`
}

type UserAuthCredential struct {
	Kind      string       `tfsdk:"kind"`
	Hash      types.String `tfsdk:"hash"`
	Email     types.String `tfsdk:"email"`
	Provider  types.String `tfsdk:"provider"`
	TotpKey   types.List   `tfsdk:"totp_key"`   //[]uint8
	PublicKey types.String `tfsdk:"public_key"` //string
}
