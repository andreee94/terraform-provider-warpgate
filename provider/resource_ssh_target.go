package provider

import (
	"context"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = sshTargetResourceType{}
var _ resource.Resource = sshTargetResource{}
var _ resource.ResourceWithImportState = sshTargetResource{}

type sshTargetResourceType struct{}

func (t sshTargetResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"allow_roles": {
				// Type: schema.TypeList,
				// Elem: &schema.Schema{
				// 	Type: schema.TypeString,
				// },
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
				Required: false,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Id of the ssh target in warpgate",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				Type:     types.StringType,
				Computed: false,
				Required: true,
			},
			"options": {
				Computed: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"host": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
					},
					"port": {
						Type:     types.Int64Type,
						Computed: false,
						Required: true,
					},
					"username": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
					},
				}),
			},
		},
	}, nil
}

func (t sshTargetResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return sshTargetResource{
		provider: provider,
	}, diags
}

type sshTargetResource struct {
	provider warpgateProvider
}

func (r sshTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.TargetSsh

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.provider.client.CreateTargetWithResponse(ctx, warpgate.CreateTargetJSONBody{
		Name: resourceState.Name,
		Options: warpgate.TargetOptionsTargetSSHOptions{
			Kind:     "Ssh",
			Host:     resourceState.Options.Host,
			Port:     resourceState.Options.Port,
			Username: resourceState.Options.Username,
			Auth: warpgate.SSHTargetAuthSshTargetPublicKeyAuth{
				Kind: "PublicKey",
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create ssh target",
			"Failed to create ssh target",
		)
		return
	}

	if response.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed to create ssh target, wrong error code.",
			fmt.Sprintf("Failed to create ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	resourceState.Id = types.String{Value: response.JSON201.Id.String()}
	resourceState.AllowRoles = response.JSON201.AllowRoles

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r sshTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.TargetSsh

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id_as_uuid, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the id as uuid",
			fmt.Sprintf("Failed to parse the id %s as uuid", resourceState.Id.String()),
		)
		return
	}

	response, err := r.provider.client.GetTargetWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read ssh target",
			fmt.Sprintf("Failed to read ssh target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read ssh target, wrong error code.",
			fmt.Sprintf("Failed to read ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	var sshoptions warpgate.TargetOptionsTargetSSHOptions
	err = mapstructure.Decode(response.JSON200.Options, &sshoptions)

	if err != nil || sshoptions.Kind != "Ssh" {
		resp.Diagnostics.AddError(
			"Failed to read ssh target. Wrong options",
			fmt.Sprintf("Failed to read ssh target %v. Wrong options type. ", response.JSON200),
		)
		return
	}

	resourceState.AllowRoles = response.JSON200.AllowRoles
	resourceState.Name = response.JSON200.Name
	resourceState.Options.Host = sshoptions.Host
	resourceState.Options.Port = sshoptions.Port
	resourceState.Options.Username = sshoptions.Username

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r sshTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceState provider_models.TargetSsh

	diags := req.Plan.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id_as_uuid, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the id as uuid",
			fmt.Sprintf("Failed to parse the id '%s' as uuid", resourceState.Id),
		)
		return
	}

	response, err := r.provider.client.UpdateTargetWithResponse(ctx, id_as_uuid, warpgate.UpdateTargetJSONBody{
		Name: resourceState.Name,
		Options: warpgate.TargetOptionsTargetSSHOptions{
			Kind:     "Ssh",
			Host:     resourceState.Options.Host,
			Port:     resourceState.Options.Port,
			Username: resourceState.Options.Username,
			Auth: warpgate.SSHTargetAuthSshTargetPublicKeyAuth{
				Kind: "PublicKey",
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update ssh target",
			fmt.Sprintf("Failed to update ssh target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to update ssh target, wrong error code.",
			fmt.Sprintf("Failed to update ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	// probably unnecessary check
	if response.JSON200.Id != id_as_uuid || response.JSON200.Name != resourceState.Name {
		resp.Diagnostics.AddWarning(
			"Created resource is different from requested.",
			fmt.Sprintf("Created resource is different from requested. Requested: (%s, %s), Created: (%s, %s)",
				response.JSON200.Id, response.JSON200.Name,
				resourceState.Id, resourceState.Name,
			),
		)
		return
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r sshTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.TargetSsh

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id_as_uuid, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the id as uuid",
			fmt.Sprintf("Failed to parse the id '%s' as uuid", resourceState.Id),
		)
		return
	}

	response, err := r.provider.client.DeleteTargetWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete ssh target",
			fmt.Sprintf("Failed to delete ssh target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to delete ssh target, wrong error code.",
			fmt.Sprintf("Failed to delete ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}
}

func (r sshTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
