package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TargetRoles struct {
	Id      types.String `tfsdk:"id"`
	RoleIds types.Set    `tfsdk:"role_ids"`
}
