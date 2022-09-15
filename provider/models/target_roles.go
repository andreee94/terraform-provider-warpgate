package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TargetRoles struct {
	// ID have to be nullable
	Id      string    `tfsdk:"id"`
	RoleIds types.Set `tfsdk:"role_ids"`
}
