package provider

import (
	"context"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/provider/validators"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
// var _ provider.ResourceType = sshTargetResourceType{}
var _ resource.Resource = &sshTargetResource{}
var _ resource.ResourceWithImportState = &sshTargetResource{}

// type sshTargetResourceType struct{}

func (r *sshTargetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"allow_roles": {
				Type:     types.SetType{ElemType: types.StringType},
				Computed: true,
				Required: false,
				Optional: true,
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
				Computed: false,
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"host": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							validators.IsDomain(),
						},
					},
					"port": {
						Type:     types.Int64Type,
						Computed: false,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							int64validator.Between(1, 65535),
						},
					},
					"username": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
					},
					"auth_kind": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.OneOf(
								string(warpgate.Password),
								string(warpgate.PublicKey),
							),
						},
					},
					"password": {
						Type:      types.StringType,
						Computed:  false,
						Required:  false,
						Optional:  true,
						Sensitive: true,
					},
				}),
			},
		},
	}, nil
}

// func (t sshTargetResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

// 	return sshTargetResource{
// 		provider: provider,
// 	}, diags
// }

func NewSshTargetResource() resource.Resource {
	return &sshTargetResource{}
}

func (r *sshTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_target"
}

func (r *sshTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type sshTargetResource struct {
	provider *warpgateProvider
}

func (r *sshTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.TargetSsh

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.provider.client.CreateTargetWithResponse(ctx, warpgate.CreateTargetJSONBody{
		Name: resourceState.Name.Value,
		Options: warpgate.TargetOptionsTargetSSHOptions{
			Kind:     "Ssh",
			Host:     resourceState.Options.Host.Value,
			Port:     uint16(resourceState.Options.Port.Value),
			Username: resourceState.Options.Username.Value,
			Auth:     GenerateSshAuth(resourceState),
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
	resourceState.AllowRoles = ArrayOfStringToTerraformSet(response.JSON201.AllowRoles)

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *sshTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	if response.StatusCode() == 404 {
		resp.Diagnostics.AddWarning(
			"Failed to read ssh target, resource not found. Removing from the state.",
			fmt.Sprintf("Failed to read ssh target. (Error code: %d)", response.StatusCode()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read ssh target, wrong error code.",
			fmt.Sprintf("Failed to read ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	sshoptions, err := ParseSshOptions(response.JSON200.Options)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read ssh target. Wrong options",
			fmt.Sprintf("Failed to read ssh target %v. Wrong options type. (Error: %v ", response.JSON200, err),
		)
		return
	}

	if sshoptions == nil {
		resp.Diagnostics.AddError(
			"Failed to read ssh target. Not an ssh target",
			"Failed to read ssh target. Not an ssh target",
		)
		return
	}

	resourceState.AllowRoles = ArrayOfStringToTerraformSet(response.JSON200.AllowRoles)
	resourceState.Name = types.String{Value: response.JSON200.Name}
	resourceState.Options = &provider_models.TargetSSHOptions{
		Host:     sshoptions.Host,
		Port:     sshoptions.Port,
		Username: sshoptions.Username,
		AuthKind: sshoptions.AuthKind,
		Password: sshoptions.Password,
	}
	// resourceState.Options.Host = sshoptions.Host
	// resourceState.Options.Port = sshoptions.Port
	// resourceState.Options.Username = sshoptions.Username
	// resourceState.Options.AuthKind = sshoptions.AuthKind
	// if resourceState.Options.AuthKind == string(warpgate.Password) {
	// 	resourceState.Options.Password = sshoptions.Password
	// }

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *sshTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourcePlan provider_models.TargetSsh

	diags := req.Plan.Get(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id_as_uuid, err := uuid.Parse(resourcePlan.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the id as uuid",
			fmt.Sprintf("Failed to parse the id '%s' as uuid", resourcePlan.Id),
		)
		return
	}

	response, err := r.provider.client.UpdateTargetWithResponse(ctx, id_as_uuid, warpgate.UpdateTargetJSONBody{
		Name: resourcePlan.Name.Value,
		Options: warpgate.TargetOptionsTargetSSHOptions{
			Kind:     "Ssh",
			Host:     resourcePlan.Options.Host.Value,
			Port:     uint16(resourcePlan.Options.Port.Value),
			Username: resourcePlan.Options.Username.Value,
			Auth:     GenerateSshAuth(resourcePlan),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update ssh target",
			fmt.Sprintf("Failed to update ssh target with id '%s'. (Error: %s)", resourcePlan.Id, err),
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
	if response.JSON200.Id != id_as_uuid || response.JSON200.Name != resourcePlan.Name.Value {
		resp.Diagnostics.AddWarning(
			"Created resource is different from requested.",
			fmt.Sprintf("Created resource is different from requested. Requested: (%s, %s), Created: (%s, %s)",
				response.JSON200.Id, response.JSON200.Name,
				resourcePlan.Id, resourcePlan.Name,
			),
		)
	}
	resourcePlan.AllowRoles = ArrayOfStringToTerraformSet(response.JSON200.AllowRoles)

	tflog.Debug(ctx, fmt.Sprintf("Updating ssh_target state: %v", resourcePlan))

	diags = resp.State.Set(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)
}

func (r *sshTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	if response.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Failed to delete ssh target, wrong error code.",
			fmt.Sprintf("Failed to delete ssh target. (Error code: %d)", response.StatusCode()),
		)
		return
	}
}

func (r *sshTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GenerateSshAuth(resourceState provider_models.TargetSsh) warpgate.SSHTargetAuth {
	var auth warpgate.SSHTargetAuth

	if resourceState.Options.AuthKind.Value == string(warpgate.Password) {
		auth = warpgate.SSHTargetAuthSshTargetPasswordAuth{
			Password: resourceState.Options.Password.Value,
			Kind:     resourceState.Options.AuthKind.Value,
		}
	} else if resourceState.Options.AuthKind.Value == string(warpgate.PublicKey) {
		auth = warpgate.SSHTargetAuthSshTargetPublicKeyAuth{
			Kind: resourceState.Options.AuthKind.Value,
		}
	}
	return auth
}

func ParseSshOptions(options warpgate.TargetOptions) (*provider_models.TargetSSHOptions, error) {
	var result provider_models.TargetSSHOptions
	var sshoptions warpgate.TargetOptionsTargetSSHOptions
	err := mapstructure.Decode(options, &sshoptions)

	if err != nil {
		return nil, err
	}

	if sshoptions.Kind != "Ssh" {
		return nil, nil
	}

	var kind struct {
		Kind string `json:"kind"`
	}
	err = mapstructure.Decode(sshoptions.Auth, &kind)

	if err != nil {
		return nil, err
	}

	result.AuthKind = types.String{Value: kind.Kind}
	result.Host = types.String{Value: sshoptions.Host}
	result.Port = types.Int64{Value: int64(sshoptions.Port)}
	result.Username = types.String{Value: sshoptions.Username}
	result.Password = types.String{Null: true}

	if kind.Kind == string(warpgate.Password) {
		var auth warpgate.SshTargetPasswordAuth

		err = mapstructure.Decode(sshoptions.Auth, &auth)

		if err != nil {
			return nil, err
		}

		result.Password = types.String{Value: auth.Password}
	}

	return &result, err
}
