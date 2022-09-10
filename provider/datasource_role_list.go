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
var _ provider.DataSourceType = roleListDataSourceType{}
var _ datasource.DataSource = roleListDataSource{}

type roleListDataSourceType struct{}

func (t roleListDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"roles": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
				}),
			},
		},
	}, nil
}

func (t roleListDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return roleListDataSource{
		provider: provider,
	}, diags
}

type roleListDataSource struct {
	provider warpgateProvider
}

func (d roleListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var resourceState struct {
		Roles []provider_models.Role `tfsdk:"roles"`
	}

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.provider.client.GetTargetsWithResponse(ctx)

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
			Id:   types.String{Value: role.Id.String()},
			Name: role.Name,
		})
	}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
