package provider

import (
	"context"
	"fmt"
	"terraform-provider-warpgate/provider/models"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &userTargetResource{}
var _ resource.ResourceWithImportState = &userTargetResource{}

func (r *userTargetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Id of the user in warpgate",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"roles": {
				Type:     types.SetType{ElemType: types.StringType},
				Computed: true,
			},
			"username": {
				Type:     types.StringType,
				Computed: false,
				Required: true,
			},
			"credentials": {
				Computed: false,
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"kind": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.OneOf(
								string(warpgate.Sso),
								string(warpgate.Totp),
								string(warpgate.Password),
								string(warpgate.PublicKey),
							),
						},
					},
					/////////////////////////////////////////////////////////////////////////////////
					"hash": {
						Type:        types.StringType,
						Computed:    false,
						Required:    false,
						Optional:    true,
						Sensitive:   true,
						Description: "Only for kind: Password",
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("email"),
								path.MatchRelative().AtParent().AtName("provider"),
								path.MatchRelative().AtParent().AtName("public_key"),
								path.MatchRelative().AtParent().AtName("totp_key"),
							),
						},
					},
					/////////////////////////////////////////////////////////////////////////////////
					"email": {
						Type:        types.StringType,
						Computed:    false,
						Required:    false,
						Optional:    true,
						Description: "Only for kind: Sso",
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("hash"),
								path.MatchRelative().AtParent().AtName("public_key"),
								path.MatchRelative().AtParent().AtName("totp_key"),
							),
							schemavalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("provider"),
							),
						},
					},
					"provider": {
						Type:        types.StringType,
						Computed:    false,
						Required:    false,
						Optional:    true,
						Description: "Only for kind: Sso",
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("hash"),
								path.MatchRelative().AtParent().AtName("public_key"),
								path.MatchRelative().AtParent().AtName("totp_key"),
							),
							schemavalidator.AlsoRequires(
								path.MatchRelative().AtParent().AtName("email"),
							),
						},
					},
					/////////////////////////////////////////////////////////////////////////////////
					"public_key": {
						Type:        types.StringType,
						Computed:    false,
						Required:    false,
						Optional:    true,
						Description: "Only for kind: PublicKey",
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("hash"),
								path.MatchRelative().AtParent().AtName("email"),
								path.MatchRelative().AtParent().AtName("provider"),
								path.MatchRelative().AtParent().AtName("totp_key"),
							),
						},
					},
					/////////////////////////////////////////////////////////////////////////////////
					"totp_key": {
						Type:        types.ListType{ElemType: types.Int64Type},
						Computed:    false,
						Required:    false,
						Optional:    true,
						Sensitive:   true,
						Description: "Only for kind: Totp",
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("hash"),
								path.MatchRelative().AtParent().AtName("email"),
								path.MatchRelative().AtParent().AtName("provider"),
								path.MatchRelative().AtParent().AtName("public_key"),
							),
						},
					},
					/////////////////////////////////////////////////////////////////////////////////
				}),
			},
		},
	}, nil
}

func NewUserResource() resource.Resource {
	return &userTargetResource{}
}

func (r *userTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type userTargetResource struct {
	provider *warpgateProvider
}

func (r *userTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.User

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.provider.client.CreateUserWithResponse(ctx, warpgate.UserDataRequest{
		Username:    resourceState.Username,
		Credentials: GenerateWarpgateUserAuthCredentials(resourceState),
		// CredentialPolicy: &warpgate.UserRequireCredentialsPolicy{},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create user",
			"Failed to create user",
		)
		return
	}

	if response.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed to create user, wrong error code.",
			fmt.Sprintf("Failed to create user. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	resourceState.Id = types.String{Value: response.JSON201.Id.String()}
	resourceState.Roles = ArrayOfStringToTerraformSet(response.JSON201.Roles)

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *userTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.User

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

	response, err := r.provider.client.GetUserWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read user",
			fmt.Sprintf("Failed to read user with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() == 404 {
		resp.Diagnostics.AddWarning(
			"Failed to read user, resource not found. Removing from the state.",
			fmt.Sprintf("Failed to read user. (Error code: %d)", response.StatusCode()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read user, wrong error code.",
			fmt.Sprintf("Failed to read user. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	user, err := ParseUser(response.JSON200)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read user. Wrong options",
			fmt.Sprintf("Failed to read user %v. Wrong options type. (Error: %v ", response.JSON200, err),
		)
		return
	}

	if user == nil {
		resp.Diagnostics.AddError(
			"Failed to read user. Not an user",
			"Failed to read user. Not an user",
		)
		return
	}

	resourceState.Roles = user.Roles
	resourceState.Credentials = user.Credentials
	resourceState.Username = user.Username

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *userTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourcePlan provider_models.User

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

	response, err := r.provider.client.UpdateUserWithResponse(ctx, id_as_uuid, warpgate.UserDataRequest{
		Username:    resourcePlan.Username,
		Credentials: GenerateWarpgateUserAuthCredentials(resourcePlan),
		// CredentialPolicy: &warpgate.UserRequireCredentialsPolicy{},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update user",
			fmt.Sprintf("Failed to update user with id '%s'. (Error: %s)", resourcePlan.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to update user, wrong error code.",
			fmt.Sprintf("Failed to update user. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	// probably unnecessary check
	if response.JSON200.Id != id_as_uuid || response.JSON200.Username != resourcePlan.Username {
		resp.Diagnostics.AddWarning(
			"Created resource is different from requested.",
			fmt.Sprintf("Created resource is different from requested. Requested: (%s, %s), Created: (%s, %s)",
				response.JSON200.Id, response.JSON200.Username,
				resourcePlan.Id, resourcePlan.Username,
			),
		)
	}
	resourcePlan.Roles = ArrayOfStringToTerraformSet(response.JSON200.Roles)

	tflog.Debug(ctx, fmt.Sprintf("Updating user state: %v", resourcePlan))

	diags = resp.State.Set(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)
}

func (r *userTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.User

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

	response, err := r.provider.client.DeleteUserWithResponse(ctx, id_as_uuid)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete user",
			fmt.Sprintf("Failed to delete user with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Failed to delete user, wrong error code.",
			fmt.Sprintf("Failed to delete user. (Error code: %d)", response.StatusCode()),
		)
		return
	}
}

func (r *userTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func ParseUserCredential(credential warpgate.UserAuthCredential) (result *models.UserAuthCredential, err error) {

	var kind struct {
		Kind string `json:"kind"`
	}
	err = mapstructure.Decode(credential, &kind)

	if err != nil {
		return nil, err
	}

	result = &models.UserAuthCredential{
		Kind:      kind.Kind,
		Hash:      types.String{Null: true},
		Email:     types.String{Null: true},
		Provider:  types.String{Null: true},
		TotpKey:   types.List{Null: true, ElemType: types.Int64Type},
		PublicKey: types.String{Null: true},
	}

	if kind.Kind == string(warpgate.Sso) {

		c := &warpgate.UserAuthCredentialUserSsoCredential{}
		err = mapstructure.Decode(credential, &c)

		var provider types.String

		if c.Provider == nil {
			provider = types.String{Null: true}
		} else {
			provider = types.String{Value: *c.Provider}
		}

		result.Email = types.String{Value: c.Email}
		result.Provider = provider

	} else if kind.Kind == string(warpgate.Password) {

		c := &warpgate.UserAuthCredentialUserPasswordCredential{}
		err = mapstructure.Decode(credential, &c)

		result.Hash = types.String{Value: c.Hash}

	} else if kind.Kind == string(warpgate.PublicKey) {

		c := &warpgate.UserAuthCredentialUserPublicKeyCredential{}
		err = mapstructure.Decode(credential, &c)

		result.PublicKey = types.String{Value: c.Key}

	} else if kind.Kind == string(warpgate.Totp) {

		c := &warpgate.UserAuthCredentialUserTotpCredential{}
		err = mapstructure.Decode(credential, &c)

		result.TotpKey = ArrayOfUint8ToTerraformList(c.Key)
	}

	return
}

func ParseUser(user *warpgate.User) (result *models.User, err error) {
	result = &provider_models.User{
		Id:       types.String{Value: user.Id.String()},
		Username: user.Username,
		Roles:    ArrayOfStringToTerraformSet(user.Roles),
	}

	for _, c := range user.Credentials {
		credential, err := ParseUserCredential(c)

		if err != nil {
			return nil, err
		}

		result.Credentials = append(result.Credentials, *credential)
	}

	return
}

func GenerateWarpgateUserAuthCredentials(user models.User) (result []warpgate.UserAuthCredential) {
	for _, c := range user.Credentials {

		var credential warpgate.UserAuthCredential

		if c.Kind == string(warpgate.Sso) {
			credential = warpgate.UserAuthCredentialUserSsoCredential{
				Kind:     c.Kind,
				Email:    c.Email.Value,
				Provider: &c.Provider.Value, // TODO check for null
			}
		} else if c.Kind == string(warpgate.Password) {
			credential = warpgate.UserAuthCredentialUserPasswordCredential{
				Kind: c.Kind,
				Hash: c.Hash.Value,
			}
		} else if c.Kind == string(warpgate.Password) {
			credential = warpgate.UserAuthCredentialUserPublicKeyCredential{
				Kind: c.Kind,
				Key:  c.PublicKey.Value,
			}
		} else if c.Kind == string(warpgate.Totp) {
			credential = warpgate.UserAuthCredentialUserTotpCredential{
				Kind: c.Kind,
				Key:  TerraformListToArrayOfUint8(c.TotpKey),
			}
		}

		result = append(result, credential)
	}

	return
}
