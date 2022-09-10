package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"terraform-provider-warpgate/provider/modifiers"
	"terraform-provider-warpgate/provider/validators"
	"terraform-provider-warpgate/warpgate"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// var stderr = os.Stderr

//	func New(version string) func() provider.Provider {
//		return func() provider.Provider {
//			return &warpgateProvider{
//				version: version,
//			}
//		}
//	}
func New() func() provider.Provider {
	return func() provider.Provider {
		return &warpgateProvider{}
	}
}

func convertProviderType(in provider.Provider) (warpgateProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*warpgateProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return warpgateProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return warpgateProvider{}, diags
	}

	return *p, diags
}

type warpgateProvider struct {
	configured bool
	// version    string
	client *warpgate.WarpgateClient
	// router     *warpgate.warpgateRouter
}

func (p *warpgateProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:        types.StringType,
				Description: "The hostname of the warpgate server",
				Optional:    false,
				Required:    true,
				Computed:    false,
				Validators: []tfsdk.AttributeValidator{
					validators.StringRegex{Regex: regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3})$|^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`)},
				},
			},
			"port": {
				Type:        types.Int64Type,
				Description: "The port of the warpgate server (Default: 8888)",
				Optional:    true,
				Required:    false,
				Validators: []tfsdk.AttributeValidator{
					validators.IntBetween(1, 65535),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.IntDefault(8888),
				},
			},
			"username": {
				Type:        types.StringType,
				Description: "The username to login to the warpgate server",
				Optional:    false,
				Computed:    false,
				Required:    true,
			},
			"password": {
				Type:        types.StringType,
				Description: "The password to login to the warpgate server",
				Optional:    false,
				Computed:    false,
				Sensitive:   true,
				Required:    true,
			},
			"insecure_skip_verify": {
				Type:        types.BoolType,
				Description: "If to skip the verification of the tls certificate",
				Optional:    true,
				Computed:    false,
				Required:    false,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.BoolDefault(false),
				},
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Host               types.String `tfsdk:"host"`
	Port               types.Int64  `tfsdk:"port"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

func (p *warpgateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var host string
	var port int
	var username string
	var password string
	var insecureSkipVerify bool

	if !checkForUnknowsInConfig(&config, resp) {
		return
	}

	if config.Host.Null {
		host = os.Getenv("WARPGATE_HOST")
	} else {
		host = config.Host.Value
	}

	if config.Username.Null {
		username = os.Getenv("WARPGATE_USERNAME")
	} else {
		username = config.Username.Value
	}

	if config.Port.Null {
		portString := os.Getenv("WARPGATE_PORT")
		if portString == "" {
			port = 8888
		} else {
			port, err = strconv.Atoi(portString)
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid port",
					"The port must be an integer",
				)
				return
			}
		}
	} else {
		port = int(config.Port.Value)
	}

	if config.Password.Null {
		password = os.Getenv("WARPGATE_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if config.InsecureSkipVerify.Null {
		envValue := os.Getenv("WARPGATE_INSECURE_SKIP_VERIFY")

		if len(envValue) > 0 {
			insecureSkipVerify, err = strconv.ParseBool(envValue)
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid insecureSkipVerify",
					"The insecureSkipVerify must be a valid bool (Valid true values: '1', 't', 'T', 'true', 'TRUE', 'True'. Valid false values: '0', 'f', 'F', 'false', 'FALSE', 'False')",
				)
				return
			}
		} else {
			insecureSkipVerify = false
		}

	} else {
		insecureSkipVerify = config.InsecureSkipVerify.Value
	}

	if username == "" {
		resp.Diagnostics.AddError(
			"Unable to find username",
			"Username cannot be an empty string",
		)
		return
	}

	p.client = warpgate.NewWarpgateClient(host, port, insecureSkipVerify)

	err = p.client.Login(username, password)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to login",
			err.Error(),
		)
		return
	}

	p.configured = true
}

func checkForUnknowsInConfig(config *providerData, resp *provider.ConfigureResponse) bool {
	if config.Host.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return false
	}

	if config.Username.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return false
	}

	if config.Password.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return false
	}
	return true
}

func (p *warpgateProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"warpgate_role":       roleResourceType{},
		"warpgate_ssh_target": sshTargetResourceType{},
		// "warpgate_port_forwarded": resourcePortForwardedType{},
	}, nil
}

func (p *warpgateProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"warpgate_ssh_target_list":  sshTargetListDataSourceType{},
		"warpgate_http_target_list": httpTargetListDataSourceType{},
		"warpgate_role_list":        roleListDataSourceType{},
		"warpgate_sshkey_list":      sshkeyListDataSourceType{},
	}, nil
}
