package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	provider_models "terraform-provider-warpgate/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &sshkeyListDataSource{}

func (d sshkeyListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"kind": schema.StringAttribute{Computed: false, Required: true},
			"sshkeys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind":              schema.StringAttribute{Computed: true},
						"public_key_base64": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func NewSshkeyListDataSource() datasource.DataSource {
	return &sshkeyListDataSource{}
}

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
			Kind:            types.StringValue(sshkey.Kind),
			PublicKeyBase64: types.StringValue(sshkey.PublicKeyBase64),
		})
	}

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.StringValue(randomUUID.String())

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
