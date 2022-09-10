package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type Role struct {
	// ID have to be nullable
	Id   types.String `tfsdk:"id"`
	Name string       `tfsdk:"name"`
}
