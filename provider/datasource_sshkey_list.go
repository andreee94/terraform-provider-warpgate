package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	provider_models "terraform-provider-warpgate/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = sshkeyListDataSourceType{}
var _ datasource.DataSource = sshkeyListDataSource{}

type sshkeyListDataSourceType struct{}

func (t sshkeyListDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
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
						Required: false,
					},
					"public_key_base64": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
				}),
			},
		},
	}, nil
}

func (t sshkeyListDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return sshkeyListDataSource{
		provider: provider,
	}, diags
}

type sshkeyListDataSource struct {
	provider warpgateProvider
}

func (d sshkeyListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var resourceState struct {
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
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
