package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// type TargetOptions interface{}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetSsh struct {
	AllowRoles types.Set        `tfsdk:"allow_roles"`
	Id         types.String     `tfsdk:"id"`
	Name       string           `tfsdk:"name"`
	Options    TargetSSHOptions `tfsdk:"options"`
}

type TargetSSHOptions struct {
	Host     string       `tfsdk:"host"`
	Port     uint16       `tfsdk:"port"`
	Username string       `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	AuthKind string       `tfsdk:"auth_kind"`
}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetHttp struct {
	AllowRoles types.Set         `tfsdk:"allow_roles"`
	Id         types.String      `tfsdk:"id"`
	Name       string            `tfsdk:"name"`
	Options    TargetHttpOptions `tfsdk:"options"`
}

type TargetHttpOptions struct {
	ExternalHost types.String               `tfsdk:"external_host"`
	Url          string                     `tfsdk:"url"`
	Tls          TargetTls                  `tfsdk:"tls"`
	Headers      *TargetHttpOptions_Headers `tfsdk:"headers"`
}

type TargetHttpOptions_Headers struct {
	AdditionalProperties map[string]string `tfsdk:"-"`
}

type TargetTls struct {
	Mode   string `tfsdk:"mode"`
	Verify bool   `tfsdk:"verify"`
}

/////////////////////////////////////////
/////////////////////////////////////////
