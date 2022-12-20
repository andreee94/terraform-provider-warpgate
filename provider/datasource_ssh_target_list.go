package provider

import (
	"context"
	"fmt"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	provider_models "terraform-provider-warpgate/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &sshTargetListDataSource{}

func NewSshTargetListDataSource() datasource.DataSource {
	return &sshTargetListDataSource{}
}

func (d sshTargetListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"targets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"allow_roles": schema.SetAttribute{Computed: true, ElementType: types.StringType},
						"options": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"host":      schema.StringAttribute{Computed: true},
								"port":      schema.Int64Attribute{Computed: true},
								"username":  schema.StringAttribute{Computed: true},
								"auth_kind": schema.StringAttribute{Computed: true},
								"password":  schema.StringAttribute{Computed: true, Sensitive: true},
							},
						},
					},
				},
			},
		},
	}
}

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

		if sshoptions == nil {
			tflog.Debug(ctx, "Not an ssh target. Continuing.")
			continue
		}

		resourceState.Targets = append(resourceState.Targets, provider_models.TargetSsh{
			// AllowRoles: target.AllowRoles,
			Id:         types.StringValue(target.Id.String()),
			Name:       types.StringValue(target.Name),
			AllowRoles: ArrayOfStringToTerraformSet(target.AllowRoles),
			Options: &provider_models.TargetSSHOptions{
				Host:     sshoptions.Host,
				Port:     sshoptions.Port,
				Username: sshoptions.Username,
				AuthKind: sshoptions.AuthKind,
				Password: If(
					sshoptions.AuthKind.ValueString() == string(warpgate.Password),
					sshoptions.Password,
					types.StringNull(),
				),
			},
		})
	}

	randomUUID, _ := uuid.NewRandom()
	resourceState.Id = types.StringValue(randomUUID.String())

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}
