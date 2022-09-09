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
var _ provider.DataSourceType = httpTargetListDataSourceType{}
var _ datasource.DataSource = httpTargetListDataSource{}

type httpTargetListDataSourceType struct{}

func (t httpTargetListDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"targets": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"allow_roles": {
						Type:     types.ListType{ElemType: types.StringType},
						Computed: true,
						Required: false,
					},
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
							"external_host": {
								Type:     types.StringType,
								Computed: true,
								Required: false,
							},
							"url": {
								Type:     types.StringType,
								Computed: true,
								Required: false,
							},
							"headers": {
								Type:     types.MapType{ElemType: types.StringType},
								Computed: true,
								Required: false,
							},
							"tls": {
								Computed: true,
								Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
									"mode": {
										Type:     types.StringType,
										Computed: true,
										Required: false,
									},
									"verify": {
										Type:     types.BoolType,
										Computed: true,
										Required: false,
									},
								}),
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (t httpTargetListDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return httpTargetListDataSource{
		provider: provider,
	}, diags
}

type httpTargetListDataSource struct {
	provider warpgateProvider
}

func (d httpTargetListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var data exampleDataSourceData

	var resourceState struct {
		Targets []provider_models.TargetHttp `tfsdk:"targets"`
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

		var httpoptions *warpgate.TargetOptionsTargetHTTPOptions
		err = mapstructure.Decode(target.Options, &httpoptions)

		if err != nil || httpoptions == nil || httpoptions.Kind != "Http" {
			tflog.Info(ctx, fmt.Sprintf("Target %v is not http, skipping.", target))
			continue
		}

		var headers *provider_models.TargetHttpOptions_Headers

		if httpoptions.Headers == nil {
			headers = nil
		} else {
			headers = &provider_models.TargetHttpOptions_Headers{
				AdditionalProperties: httpoptions.Headers.AdditionalProperties,
			}
		}

		resourceState.Targets = append(resourceState.Targets, provider_models.TargetHttp{
			AllowRoles: target.AllowRoles,
			Id:         target.Id.String(),
			Name:       target.Name,
			Options: provider_models.TargetHttpOptions{
				ExternalHost: httpoptions.ExternalHost,
				Url:          httpoptions.Url,
				Headers:      headers,
				Tls: provider_models.TargetTls{
					Mode:   string(httpoptions.Tls.Mode),
					Verify: httpoptions.Tls.Verify,
				},
			},
		})
	}
	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
