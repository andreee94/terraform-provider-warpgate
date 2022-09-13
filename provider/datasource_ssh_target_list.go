package provider

import (
	"context"
	"fmt"
	"terraform-provider-warpgate/warpgate"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/mapstructure"

	provider_models "terraform-provider-warpgate/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = sshTargetListDataSourceType{}
var _ datasource.DataSource = sshTargetListDataSource{}

type sshTargetListDataSourceType struct{}

func (t sshTargetListDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"targets": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					// "allow_roles": {
					// 	// Type: schema.TypeList,
					// 	// Elem: &schema.Schema{
					// 	// 	Type: schema.TypeString,
					// 	// },
					// 	Type:     types.ListType{ElemType: types.StringType},
					// 	Computed: true,
					// 	Required: false,
					// },
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
					"options": {
						Computed: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"host": {
								Type:     types.StringType,
								Computed: true,
								Required: false,
							},
							"port": {
								Type:     types.Int64Type,
								Computed: true,
								Required: false,
							},
							"username": {
								Type:     types.StringType,
								Computed: true,
								Required: false,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (t sshTargetListDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return sshTargetListDataSource{
		provider: provider,
	}, diags
}

type sshTargetListDataSource struct {
	provider warpgateProvider
}

func (d sshTargetListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var data exampleDataSourceData

	var resourceState struct {
		Targets []provider_models.TargetSsh `tfsdk:"targets"`
	}

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.provider.client.GetTargetsWithResponse(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get target list",
			"Failed to get target list",
		)
		return
	}

	if response.HTTPResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Failed to get target list, wrong error code.",
			fmt.Sprintf("Failed to get target list. (Error code: %d)", response.HTTPResponse.StatusCode),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d targets.", len(*response.JSON200)))

	for _, target := range *response.JSON200 {

		tflog.Trace(ctx, fmt.Sprintf("Found %v", target))

		var sshoptions warpgate.TargetOptionsTargetSSHOptions
		err = mapstructure.Decode(target.Options, &sshoptions)

		if err != nil || sshoptions.Kind != "Ssh" {
			tflog.Info(ctx, fmt.Sprintf("Target %v is not ssh, skipping.", target))
			continue
		}

		resourceState.Targets = append(resourceState.Targets, provider_models.TargetSsh{
			// AllowRoles: target.AllowRoles,
			Id:   types.String{Value: target.Id.String()},
			Name: target.Name,
			Options: provider_models.TargetSSHOptions{
				Host:     sshoptions.Host,
				Port:     sshoptions.Port,
				Username: sshoptions.Username,
			},
		})

		// resourceState.Targets = append(resourceState.Targets, provider_models.Target{
		// 	// AllowRoles: rules,
		// 	Id:   types.String{Value: target.Id.String()},
		// 	Name: types.String{Value: target.Name},
		// 	Options: provider_models.TargetOptionsTargetSSHOptions{
		// 		Host: types.String{Value: sshoptions.Host},
		// 		// Kind:     types.String{Value: sshoptions.Kind},
		// 		Port:     types.Int64{Value: int64(sshoptions.Port)},
		// 		Username: types.String{Value: sshoptions.Username},
		// 	},
		// })
	}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
