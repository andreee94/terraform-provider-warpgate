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
var _ datasource.DataSource = &roleListDataSource{}

func NewRoleListDataSource() datasource.DataSource {
	return &roleListDataSource{}
}

func (r *roleListDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": { // required for acceptance testing
				Type:     types.StringType,
				Computed: true,
			},
			"roles": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

type roleListDataSource struct {
	provider *warpgateProvider
}

func (d *roleListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_list"
}

func (d *roleListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var resourceState struct {
		Id    types.String           `tfsdk:"id"`
		Roles []provider_models.Role `tfsdk:"roles"`
	}

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.provider.client.GetRolesWithResponse(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get role list",
			"Failed to get role list",
		)
		return
	}

	if response.HTTPResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Failed to get role list, wrong error code.",
			fmt.Sprintf("Failed to get role list. (Error code: %d)", response.HTTPResponse.StatusCode),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d roles.", len(*response.JSON200)))

	for _, role := range *response.JSON200 {

		tflog.Trace(ctx, fmt.Sprintf("Found %v", role))

		resourceState.Roles = append(resourceState.Roles, provider_models.Role{
			Id:   types.StringValue(role.Id.String()),
			Name: types.StringValue(role.Name),
		})
	}

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.StringValue(randomUUID.String())

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
