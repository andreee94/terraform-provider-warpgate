package provider

import (
	"context"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/provider/validators"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
// var _ provider.ResourceType = userRolesResourceType{}
var _ resource.Resource = &userRolesResource{}
var _ resource.ResourceWithImportState = &userRolesResource{}

// type userRolesResourceType struct{}

func (r *userRolesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				MarkdownDescription: "Id of the user in warpgate",
				Computed:            false,
				Required:            true,
				Optional:            false,
				Validators: []tfsdk.AttributeValidator{
					validators.IsUUID(),
				},
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
				Validators: []tfsdk.AttributeValidator{
					setvalidator.ValuesAre(
						validators.IsUUID(),
					),
				},
			},
		},
	}, nil
}

// func (t userRolesResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

// 	return userRolesResource{
// 		provider: provider,
// 	}, diags
// }

func NewUserRolesResource() resource.Resource {
	return &userRolesResource{}
}

func (r *userRolesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_roles"
}

func (r *userRolesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type userRolesResource struct {
	provider *warpgateProvider
}

func (r *userRolesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.UserRoles

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	userUUID, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse user id.",
			fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	for _, roleId := range resourceState.RoleIds.Elems {

		roleUUID, err := uuid.Parse(roleId.String())

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse user id.",
				fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}

		response, err := r.provider.client.AddUserRoleWithResponse(ctx, userUUID, roleUUID)

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
				fmt.Sprintf("Failed to add role %s to user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *userRolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.UserRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	userUUID, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse user id.",
			fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	response, err := r.provider.client.GetUserRolesWithResponse(ctx, userUUID)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read user roles",
			fmt.Sprintf("Failed to read roles of user with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read user roles, wrong error code.",
			fmt.Sprintf("Failed to read roles of user with id '%s'. (Error code: %d)", resourceState.Id, response.StatusCode()),
		)
		return
	}

	if response.JSON200 != nil {
		resourceState.RoleIds = ArrayOfRolesToTerraformSet(*response.JSON200)
	} else {
		resourceState.RoleIds = types.Set{ElemType: types.StringType}
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *userRolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourceState provider_models.UserRoles
	var resourcePlan provider_models.UserRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	userUUID, err := uuid.Parse(resourcePlan.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse user id.",
			fmt.Sprintf("Invalid user id %s (Err: %s)", resourcePlan.Id, err),
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
				"Failed to parse user id.",
				fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.DeleteUserRoleWithResponse(ctx, userUUID, roleUUID)

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
				fmt.Sprintf("Failed to remove role %s from user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
		} else if response.StatusCode() != 204 {
			resp.Diagnostics.AddError(
				"Failed to delete role, wrong error code.",
				fmt.Sprintf("Failed to remove role %s from user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	for _, roleId := range toBeCreated {
		roleUUID, err := uuid.Parse(roleId)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse user id.",
				fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.AddUserRoleWithResponse(ctx, userUUID, roleUUID)

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
				fmt.Sprintf("Failed to add role %s to user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)
}

func (r *userRolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.UserRoles

	diags := req.State.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	userUUID, err := uuid.Parse(resourceState.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse user id.",
			fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
		)
		return
	}

	for _, roleId := range resourceState.RoleIds.Elems {

		roleUUID, err := uuid.Parse(roleId.String())

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse user id.",
				fmt.Sprintf("Invalid user id %s (Err: %s)", resourceState.Id, err),
			)
			return
		}
		response, err := r.provider.client.DeleteUserRoleWithResponse(ctx, userUUID, roleUUID)

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
				fmt.Sprintf("Failed to remove role %s from user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
		} else if response.StatusCode() != 204 {
			resp.Diagnostics.AddError(
				"Failed to delete role, wrong error code.",
				fmt.Sprintf("Failed to remove role %s from user %s. (Error code: %d)", roleId, resourceState.Id, response.StatusCode()),
			)
			return
		}
	}
}

func (r *userRolesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
