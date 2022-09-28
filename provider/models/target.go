package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// type TargetOptions interface{}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetSsh struct {
	AllowRoles types.Set         `tfsdk:"allow_roles"`
	Id         types.String      `tfsdk:"id"`
	Name       types.String      `tfsdk:"name"`
	Options    *TargetSSHOptions `tfsdk:"options"`
}

type TargetSSHOptions struct {
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"` // uint16
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	AuthKind types.String `tfsdk:"auth_kind"`
}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetHttp struct {
	AllowRoles types.Set          `tfsdk:"allow_roles"`
	Id         types.String       `tfsdk:"id"`
	Name       types.String       `tfsdk:"name"`
	Options    *TargetHttpOptions `tfsdk:"options"`
}

type TargetHttpOptions struct {
	ExternalHost types.String               `tfsdk:"external_host"`
	Url          types.String               `tfsdk:"url"`
	Tls          *TargetTls                 `tfsdk:"tls"`
	Headers      *TargetHttpOptions_Headers `tfsdk:"headers"`
}

type TargetHttpOptions_Headers struct {
	AdditionalProperties map[string]string `tfsdk:"-"`
}

type TargetTls struct {
	Mode   types.String `tfsdk:"mode"`
	Verify types.Bool   `tfsdk:"verify"`
}

/////////////////////////////////////////
/////////////////////////////////////////
