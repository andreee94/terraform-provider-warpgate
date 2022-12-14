package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"terraform-provider-warpgate/provider/validators"
	"terraform-provider-warpgate/warpgate"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &warpgateProvider{
			version: version,
		}
	}
}

var _ provider.Provider = &warpgateProvider{}

// var _ provider.ProviderWithMetaSchema = &warpgateProvider{}

type warpgateProvider struct {
	configured bool
	version    string
	client     *warpgate.WarpgateClient
}

func (p *warpgateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "warpgate"
	resp.Version = p.version
}

// Schema(context.Context, SchemaRequest, *SchemaResponse)
func (d warpgateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The hostname of the warpgate server",
				Optional:    true,
				Validators: []validator.String{
					validators.IsDomain(),
				},
			},
			"port": schema.Int64Attribute{
				Description: "The hostname of the warpgate server",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username to login to the warpgate server",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to login to the warpgate server",
				Optional:    true,
				Sensitive:   true,
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Description: "If to skip the verification of the tls certificate (For self signed certificates)",
				Optional:    true,
			},
		},
	}
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

	if config.Host.IsNull() {
		host = os.Getenv("WARPGATE_HOST")
	} else {
		host = config.Host.ValueString()
	}

	if config.Username.IsNull() {
		username = os.Getenv("WARPGATE_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	if config.Port.IsNull() {
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
		port = int(config.Port.ValueInt64())
	}

	if config.Password.IsNull() {
		password = os.Getenv("WARPGATE_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	if config.InsecureSkipVerify.IsNull() {
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
		insecureSkipVerify = config.InsecureSkipVerify.ValueBool()
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
			// "Unable to login",
			fmt.Sprintf("Unable to login, %s:%s@%s:%d", username, password, host, port),
			err.Error(),
		)
		return
	}

	p.configured = true

	resp.DataSourceData = p
	resp.ResourceData = p
}

func checkForUnknowsInConfig(config *providerData, resp *provider.ConfigureResponse) bool {
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return false
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return false
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return false
	}
	return true
}

func (p *warpgateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHttpTargetResource,
		NewSshTargetResource,
		NewRoleResource,
		NewTargetRolesResource,
		NewUserResource,
		NewUserRolesResource,
	}
}

func (p *warpgateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRoleListDataSource,
		NewSshkeyListDataSource,
		NewSshTargetListDataSource,
		NewHttpTargetListDataSource,
	}
}
