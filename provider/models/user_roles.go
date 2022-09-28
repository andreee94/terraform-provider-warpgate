package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UserRoles struct {
	// ID have to be nullable
	Id      types.String `tfsdk:"id"`
	RoleIds types.Set    `tfsdk:"role_ids"`
}
