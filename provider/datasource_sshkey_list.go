package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	provider_models "terraform-provider-warpgate/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces
// var _ provider.DataSourceType = sshkeyListDataSourceType{}
var _ datasource.DataSource = &sshkeyListDataSource{}

// type sshkeyListDataSourceType struct{}

func (d *sshkeyListDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": { // required for acceptance testing
				Type:     types.StringType,
				Computed: true,
			},
			"kind": {
				Type:     types.StringType,
				Computed: false,
				Required: true,
				Optional: false,
			},
			"sshkeys": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"kind": {
						Type:     types.StringType,
						Computed: true,
					},
					"public_key_base64": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func NewSshkeyListDataSource() datasource.DataSource {
	return &sshkeyListDataSource{}
}

// func (d *sshkeyListDataSource) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

// 	return sshkeyListDataSource{
// 		provider: provider,
// 	}, diags
// }

type sshkeyListDataSource struct {
	provider *warpgateProvider
}

func (d *sshkeyListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sshkey_list"
}

func (d *sshkeyListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*warpgateProvider)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *warpgateProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	if !provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"Expected a configured provider but it wasn't. Please report this issue to the provider developers.",
		)

		return
	}

	d.provider = provider

}
func (d *sshkeyListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var resourceState struct {
		Id      types.String             `tfsdk:"id"`
		Kind    string                   `tfsdk:"kind"`
		SshKeys []provider_models.SshKey `tfsdk:"sshkeys"`
	}

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.provider.client.GetSshOwnKeysWithResponse(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get sshkey list",
			"Failed to get sshkey list",
		)
		return
	}

	if response.HTTPResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Failed to get sshkey list, wrong error code.",
			fmt.Sprintf("Failed to get sshkey list. (Error code: %d)", response.HTTPResponse.StatusCode),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d sshkeys.", len(*response.JSON200)))

	for _, sshkey := range *response.JSON200 {

		tflog.Trace(ctx, fmt.Sprintf("Found %v", sshkey))

		if sshkey.Kind != resourceState.Kind {
			continue
		}

		resourceState.SshKeys = append(resourceState.SshKeys, provider_models.SshKey{
			Kind:            sshkey.Kind,
			PublicKeyBase64: sshkey.PublicKeyBase64,
		})
	}

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.String{Value: randomUUID.String()}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
