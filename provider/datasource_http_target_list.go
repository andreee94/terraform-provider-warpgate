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
var _ datasource.DataSource = &httpTargetListDataSource{}

func NewHttpTargetListDataSource() datasource.DataSource {
	return &httpTargetListDataSource{}
}

func (d *httpTargetListDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": { // required for acceptance testing
				Type:     types.StringType,
				Computed: true,
			},
			"targets": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"allow_roles": {
						Type:     types.SetType{ElemType: types.StringType},
						Computed: true,
					},
					"id": {
						Type:     types.StringType,
						Computed: true,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
					"options": {
						Computed: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"external_host": {
								Type:     types.StringType,
								Computed: true,
							},
							"url": {
								Type:     types.StringType,
								Computed: true,
							},
							"headers": {
								Type:     types.MapType{ElemType: types.StringType},
								Computed: true,
							},
							"tls": {
								Computed: true,
								Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
									"mode": {
										Type:     types.StringType,
										Computed: true,
									},
									"verify": {
										Type:     types.BoolType,
										Computed: true,
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

type httpTargetListDataSource struct {
	provider *warpgateProvider
}

func (d *httpTargetListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_target_list"
}

func (d *httpTargetListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *httpTargetListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var data exampleDataSourceData

	var resourceState struct {
		Id      types.String                 `tfsdk:"id"`
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

		httpoptions, err := ParseHttpOptions(target.Options)

		if err != nil || httpoptions == nil {
			tflog.Info(ctx, fmt.Sprintf("Target %v is not http, skipping.", target))
			continue
		}

		resourceState.Targets = append(resourceState.Targets, provider_models.TargetHttp{
			AllowRoles: ArrayOfStringToTerraformSet(target.AllowRoles),
			Id:         types.StringValue(target.Id.String()),
			Name:       types.StringValue(target.Name),
			Options: &provider_models.TargetHttpOptions{
				ExternalHost: httpoptions.ExternalHost,
				Url:          httpoptions.Url,
				Headers:      httpoptions.Headers,
				Tls: &provider_models.TargetTls{
					Mode:   types.StringValue(httpoptions.Tls.Mode.ValueString()),
					Verify: httpoptions.Tls.Verify,
				},
			},
		})
	}

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.StringValue(randomUUID.String())

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
