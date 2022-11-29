package models

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// type JSONableSlice []uint8

// func (u JSONableSlice) MarshalJSON() ([]byte, error) {
// 	var result string
// 	if u == nil {
// 		result = "null"
// 	} else {
// 		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
// 	}
// 	return []byte(result), nil
// }

type User struct {
	// CredentialPolicy *UserRequireCredentialsPolicy `json:"credential_policy,omitempty"`
	Id          types.String `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	Credentials types.Set    `tfsdk:"credentials"` // []UserAuthCredential
	Roles       types.Set    `tfsdk:"roles"`
}

type UserAuthCredential struct {
	// Id        types.String `tfsdk:"id"`
	Kind      types.String `tfsdk:"kind"`
	Hash      types.String `tfsdk:"hash"`
	Email     types.String `tfsdk:"email"`
	Provider  types.String `tfsdk:"provider"`
	TotpKey   types.List   `tfsdk:"totp_key"`   //[]uint8
	PublicKey types.String `tfsdk:"public_key"` //string
}

func (u User) CredentialsAsArray(ctx context.Context) ([]UserAuthCredential, error) {
	if u.Credentials.IsNull() {
		return nil, nil
	}

	var vars []UserAuthCredential
	err := u.Credentials.ElementsAs(ctx, &vars, true)
	if err != nil {
		return nil, fmt.Errorf("error reading credentials: %s", err)
	}
	return vars, nil
}
