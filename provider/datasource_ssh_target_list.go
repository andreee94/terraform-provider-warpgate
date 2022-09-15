package provider

import (
	"context"
	"fmt"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/provider/validators"
)

// Ensure provider defined types fully satisfy framework interfaces
// var _ provider.DataSourceType = sshTargetListDataSourceType{}
var _ datasource.DataSource = &sshTargetListDataSource{}

func NewSshTargetListDataSource() datasource.DataSource {
	return &sshTargetListDataSource{}
}

// type sshTargetListDataSourceType struct{}

func (r *sshTargetListDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
							"host": {
								Type:     types.StringType,
								Computed: true,
							},
							"port": {
								Type:     types.Int64Type,
								Computed: true,
							},
							"username": {
								Type:     types.StringType,
								Computed: true,
							},
							"auth_kind": {
								Type:     types.StringType,
								Computed: true,
								Validators: []tfsdk.AttributeValidator{
									validators.StringIn([]string{string(warpgate.Password), string(warpgate.PublicKey)}, false),
								},
							},
							"password": {
								Type:      types.StringType,
								Computed:  true,
								Sensitive: true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

// func (t sshTargetListDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

// 	return sshTargetListDataSource{
// 		provider: provider,
// 	}, diags
// }

type sshTargetListDataSource struct {
	provider *warpgateProvider
}

func (d *sshTargetListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_target_list"
}

func (d *sshTargetListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sshTargetListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var data exampleDataSourceData

	var resourceState struct {
		Id      types.String                `tfsdk:"id"`
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

		sshoptions, err := ParseSshOptions(target.Options)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read ssh target. Wrong options",
				fmt.Sprintf("Failed to read ssh target %v. Wrong options type. (Error: %v ", response.JSON200, err),
			)
			return
		}

		// var sshoptions warpgate.TargetOptionsTargetSSHOptions
		// err = mapstructure.Decode(target.Options, &sshoptions)

		// if err != nil || sshoptions.Kind != "Ssh" {
		// 	tflog.Info(ctx, fmt.Sprintf("Target %v is not ssh, skipping.", target))
		// 	continue
		// }

		resourceState.Targets = append(resourceState.Targets, provider_models.TargetSsh{
			// AllowRoles: target.AllowRoles,
			Id:         types.String{Value: target.Id.String()},
			Name:       target.Name,
			AllowRoles: ArrayOfStringToTerraformSet(target.AllowRoles),
			Options: provider_models.TargetSSHOptions{
				Host:     sshoptions.Host,
				Port:     sshoptions.Port,
				Username: sshoptions.Username,
				AuthKind: sshoptions.AuthKind,
				Password: If(
					sshoptions.AuthKind == string(warpgate.Password),
					sshoptions.Password,
					types.String{Null: true},
				),
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

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.String{Value: randomUUID.String()}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
