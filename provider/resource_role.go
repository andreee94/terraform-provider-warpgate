package provider

import (
	"context"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = roleResourceType{}
var _ resource.Resource = roleResource{}
var _ resource.ResourceWithImportState = roleResource{}

type roleResourceType struct{}

func (t roleResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:     types.StringType,
				Computed: false,
				Required: true,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Id of the role in warpgate",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t roleResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return roleResource{
		provider: provider,
	}, diags
}

type roleResource struct {
	provider warpgateProvider
}

func (r roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.Role

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.provider.client.CreateRoleWithResponse(ctx, warpgate.CreateRoleJSONBody{
		Name: resourceState.Name,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create role",
			"Failed to create role",
		)
		return
	}

	if response.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed to create role, wrong error code.",
			fmt.Sprintf("Failed to create role. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	resourceState.Id = types.String{Value: response.JSON201.Id.String()}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.CreateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.String{Value: "example-id"}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	// tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.Role

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

	response, err := r.provider.client.GetRoleWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read role",
			fmt.Sprintf("Failed to read role with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read role, wrong error code.",
			fmt.Sprintf("Failed to read role. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	resourceState.Name = response.JSON200.Name

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.ReadExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceState provider_models.Role

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

	response, err := r.provider.client.UpdateRoleWithResponse(ctx, id_as_uuid, warpgate.UpdateRoleJSONBody{
		Name: resourceState.Name,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update role",
			fmt.Sprintf("Failed to update role with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to update role, wrong error code.",
			fmt.Sprintf("Failed to update role. (Error code: %d)", response.StatusCode()),
		)
		return
	}

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
	resourceState.Id = types.String{Value: response.JSON200.Id.String()}
	resourceState.Name = response.JSON200.Name

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.UpdateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.Role

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

	response, err := r.provider.client.DeleteRoleWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete role",
			fmt.Sprintf("Failed to delete role with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Failed to delete role, wrong error code.",
			fmt.Sprintf("Failed to delete role. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
