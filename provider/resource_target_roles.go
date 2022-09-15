package provider

import (
	"context"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
// var _ provider.ResourceType = targetRolesResourceType{}
var _ resource.Resource = &targetRolesResource{}
var _ resource.ResourceWithImportState = &targetRolesResource{}

// type targetRolesResourceType struct{}

func (r *targetRolesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				MarkdownDescription: "Id of the target in warpgate",
				Computed:            false,
				Required:            true,
				Optional:            false,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"role_ids": {
				Type:                types.SetType{ElemType: types.StringType},
				MarkdownDescription: "List of id roles'",
				Computed:            false,
				Required:            true,
				Optional:            false,
			},
		},
	}, nil
}

// func (t targetRolesResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

// 	return targetRolesResource{
// 		provider: provider,
// 	}, diags
// }

func NewTargetRolesResource() resource.Resource {
	return &targetRolesResource{}
}

func (r *targetRolesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_roles"
}

func (r *targetRolesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.provider = provider
}

type targetRolesResource struct {
	provider *warpgateProvider
}

func (r *targetRolesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.TargetRoles

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	targetUUID, err := uuid.Parse(resourceState.Id)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse target id.",
			fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	// for _, roleId := range GetArraySortedToString(resourceState.RoleIds) {
	for _, roleId := range resourceState.RoleIds.Elems {

		roleUUID, err := uuid.Parse(roleId.String())

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse target id.",
				fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.AddTargetRoleWithResponse(ctx, targetUUID, roleUUID)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create role",
				fmt.Sprintf("Failed to create role (Err: %s)", err),
			)
			return
		}

		if response.StatusCode() != 201 {
			resp.Diagnostics.AddError(
				"Failed to create role, wrong error code.",
				fmt.Sprintf("Failed to add role %s to target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *targetRolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.TargetRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	targetUUID, err := uuid.Parse(resourceState.Id)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse target id.",
			fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	response, err := r.provider.client.GetTargetRolesWithResponse(ctx, targetUUID)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read target roles",
			fmt.Sprintf("Failed to read roles of target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read target roles, wrong error code.",
			fmt.Sprintf("Failed to read roles of target with id '%s'. (Error code: %d)", resourceState.Id, response.StatusCode()),
		)
		return
	}

	if response.JSON200 != nil {
		resourceState.RoleIds = ArrayOfRolesToTerraformSet(*response.JSON200)

		// for _, role := range *response.JSON200 {
		// 	resourceState.RoleIds = append(resourceState.RoleIds, role.Id.String())
		// }
	} else {
		resourceState.RoleIds = types.Set{ElemType: types.StringType}
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *targetRolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceState provider_models.TargetRoles
	var resourcePlan provider_models.TargetRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	targetUUID, err := uuid.Parse(resourcePlan.Id)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse target id.",
			fmt.Sprintf("Invalid target id %s (Err: %s)", resourcePlan.Id, err),
		)
		return
	}

	resourceStateRoleIds := []string{}
	resourcePlanRoleIds := []string{}

	resourceState.RoleIds.ElementsAs(ctx, &resourceStateRoleIds, true)
	resourcePlan.RoleIds.ElementsAs(ctx, &resourcePlanRoleIds, true)

	_, toBeCreated, toBeDeleted := ArrayIntersection(resourcePlanRoleIds, resourceStateRoleIds)

	for _, roleId := range toBeDeleted {
		roleUUID, err := uuid.Parse(roleId)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse target id.",
				fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.DeleteTargetRoleWithResponse(ctx, targetUUID, roleUUID)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete role",
				fmt.Sprintf("Failed to delete role (Err: %s)", err),
			)
			return
		}

		if response.StatusCode() == 409 {
			resp.Diagnostics.AddWarning(
				"Failed to delete role, conflict.",
				fmt.Sprintf("Failed to remove role %s from target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
		} else if response.StatusCode() != 204 {
			resp.Diagnostics.AddError(
				"Failed to delete role, wrong error code.",
				fmt.Sprintf("Failed to remove role %s from target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	for _, roleId := range toBeCreated {
		roleUUID, err := uuid.Parse(roleId)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse target id.",
				fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.AddTargetRoleWithResponse(ctx, targetUUID, roleUUID)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create role",
				fmt.Sprintf("Failed to create role (Err: %s)", err),
			)
			return
		}

		if response.StatusCode() != 201 {
			resp.Diagnostics.AddError(
				"Failed to create role, wrong error code.",
				fmt.Sprintf("Failed to add role %s to target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)
}

func (r *targetRolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.TargetRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	targetUUID, err := uuid.Parse(resourceState.Id)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse target id.",
			fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	for _, roleId := range resourceState.RoleIds.Elems {

		roleUUID, err := uuid.Parse(roleId.String())

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse target id.",
				fmt.Sprintf("Invalid target id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.DeleteTargetRoleWithResponse(ctx, targetUUID, roleUUID)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete role",
				fmt.Sprintf("Failed to delete role (Err: %s)", err),
			)
			return
		}

		if response.StatusCode() == 409 {
			resp.Diagnostics.AddWarning(
				"Failed to delete role, conflict.",
				fmt.Sprintf("Failed to remove role %s from target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
		} else if response.StatusCode() != 204 {
			resp.Diagnostics.AddError(
				"Failed to delete role, wrong error code.",
				fmt.Sprintf("Failed to remove role %s from target %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	// id_as_uuid, err := uuid.Parse(resourceState.Id.Value)

	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Failed to parse the id as uuid",
	// 		fmt.Sprintf("Failed to parse the id '%s' as uuid", resourceState.Id),
	// 	)
	// 	return
	// }

	// response, err := r.provider.client.DeleteRoleWithResponse(ctx, id_as_uuid)

	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Failed to delete role",
	// 		fmt.Sprintf("Failed to delete role with id '%s'. (Error: %s)", resourceState.Id, err),
	// 	)
	// 	return
	// }

	// if response.StatusCode() != 204 {
	// 	resp.Diagnostics.AddError(
	// 		"Failed to delete role, wrong error code.",
	// 		fmt.Sprintf("Failed to delete role. (Error code: %d)", response.StatusCode()),
	// 	)
	// 	return
	// }

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *targetRolesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
