package models

// type TargetOptions interface{}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetSsh struct {
	AllowRoles []string         `tfsdk:"allow_roles"`
	Id         string           `tfsdk:"id"`
	Name       string           `tfsdk:"name"`
	Options    TargetSSHOptions `tfsdk:"options"`
}

type TargetSSHOptions struct {
	Host     string `tfsdk:"host"`
	Port     int    `tfsdk:"port"`
	Username string `tfsdk:"username"`
}

/////////////////////////////////////////
/////////////////////////////////////////

type TargetHttp struct {
	AllowRoles []string          `tfsdk:"allow_roles"`
	Id         string            `tfsdk:"id"`
	Name       string            `tfsdk:"name"`
	Options    TargetHttpOptions `tfsdk:"options"`
}

type TargetHttpOptions struct {
	ExternalHost *string                    `tfsdk:"external_host"`
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