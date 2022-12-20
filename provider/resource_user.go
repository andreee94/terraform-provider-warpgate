package provider

import (
	"context"
	"errors"
	"fmt"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &userTargetResource{}
var _ resource.ResourceWithImportState = &userTargetResource{}

var credentialsAttributes = map[string]attr.Type{
	"kind":       types.StringType,
	"hash":       types.StringType,
	"email":      types.StringType,
	"provider":   types.StringType,
	"public_key": types.StringType,
	"totp_key":   types.ListType{ElemType: types.Int64Type},
}

func (r userTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the user in warpgate",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "The list of roles that the user belong to. To assign a new role to the user refer to [user_roles](user_roles.md) ",
				ElementType:         types.StringType,
			},
			"username": schema.StringAttribute{
				Computed:            false,
				Required:            true,
				MarkdownDescription: "The username of the user.",
			},
			"credentials": schema.SetNestedAttribute{
				Computed: false,
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							Computed: false,
							Required: true,
							MarkdownDescription: "The credential type. Valid values are:\n" +
								"	- `Sso` requires: `email` and `provider`.\n" +
								"	- `Totp` requires: `totp_key`.\n" +
								"	- `Password` requires: `hash`.\n" +
								"	- `PublicKey` requires: `public_key`.\n",
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(warpgate.Sso),
									string(warpgate.Totp),
									string(warpgate.Password),
									string(warpgate.PublicKey),
								),
							},
						},
						/////////////////////////////////////////////////////////////////////////////////
						"hash": schema.StringAttribute{
							Computed:            false,
							Required:            false,
							Optional:            true,
							Sensitive:           true, // TODO it's really sensitive?
							MarkdownDescription: "The hashed password. Only for kind: `Password`",
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("email"),
									path.MatchRelative().AtParent().AtName("provider"),
									path.MatchRelative().AtParent().AtName("public_key"),
									path.MatchRelative().AtParent().AtName("totp_key"),
								),
							},
						},
						/////////////////////////////////////////////////////////////////////////////////
						"email": schema.StringAttribute{
							Computed:            false,
							Required:            false,
							Optional:            true,
							MarkdownDescription: "The email of the user in the sso system. Only for kind: `Sso`",
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("hash"),
									path.MatchRelative().AtParent().AtName("public_key"),
									path.MatchRelative().AtParent().AtName("totp_key"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("provider"),
								),
							},
						},
						"provider": schema.StringAttribute{
							Computed:            false,
							Required:            false,
							Optional:            true,
							MarkdownDescription: "The sso provider name defined in the configuration file. Only for kind: `Sso`",
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("hash"),
									path.MatchRelative().AtParent().AtName("public_key"),
									path.MatchRelative().AtParent().AtName("totp_key"),
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("email"),
								),
							},
						},
						/////////////////////////////////////////////////////////////////////////////////
						"public_key": schema.StringAttribute{
							Computed:    false,
							Required:    false,
							Optional:    true,
							Description: "The ssh public key that the user uses to connect via ssh. Only for kind: `PublicKey`",
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("hash"),
									path.MatchRelative().AtParent().AtName("email"),
									path.MatchRelative().AtParent().AtName("provider"),
									path.MatchRelative().AtParent().AtName("totp_key"),
								),
							},
						},
						/////////////////////////////////////////////////////////////////////////////////
						"totp_key": schema.ListAttribute{
							ElementType: types.Int64Type,
							Computed:    false,
							Required:    false,
							Optional:    true,
							Sensitive:   true,
							Description: "The totp secret key as array of uint8. Only for kind: `Totp`",
							Validators: []validator.List{
								listvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("hash"),
									path.MatchRelative().AtParent().AtName("email"),
									path.MatchRelative().AtParent().AtName("provider"),
									path.MatchRelative().AtParent().AtName("public_key"),
								),
							},
						},
						/////////////////////////////////////////////////////////////////////////////////

					},
				},
				// Validators: []validator.List{
				// 	listvalidator.SizeAtMost(2),
				// },
			},
		},
	}
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
		Username:    resourceState.Username.ValueString(),
		Credentials: GenerateWarpgateUserAuthCredentials(ctx, resourceState),
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

	resourceState.Id = types.StringValue(response.JSON201.Id.String())
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

	id_as_uuid, err := uuid.Parse(resourceState.Id.ValueString())

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
	// resourceState.Credentials = types.SetNull(types.ObjectType{AttrTypes: credentialsAttributes})
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

	id_as_uuid, err := uuid.Parse(resourcePlan.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the id as uuid",
			fmt.Sprintf("Failed to parse the id '%s' as uuid", resourcePlan.Id),
		)
		return
	}

	response, err := r.provider.client.UpdateUserWithResponse(ctx, id_as_uuid, warpgate.UserDataRequest{
		Username:    resourcePlan.Username.ValueString(),
		Credentials: GenerateWarpgateUserAuthCredentials(ctx, resourcePlan),
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
	if response.JSON200.Id != id_as_uuid || response.JSON200.Username != resourcePlan.Username.ValueString() {
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

	id_as_uuid, err := uuid.Parse(resourceState.Id.ValueString())

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

func ParseUserCredential(credential warpgate.UserAuthCredential) (result types.Object, err error) {

	discriminator, err := credential.Discriminator()
	if err != nil {
		return types.ObjectNull(credentialsAttributes), err
	}

	value := map[string]attr.Value{
		"kind":       types.StringNull(),
		"hash":       types.StringNull(),
		"email":      types.StringNull(),
		"provider":   types.StringNull(),
		"public_key": types.StringNull(),
		"totp_key":   types.ListNull(types.Int64Type),
	}

	switch discriminator {
	case "Password":
		auth, err := credential.AsUserAuthCredentialUserPasswordCredential()

		if err != nil {
			return types.ObjectNull(credentialsAttributes), err
		}

		value["kind"] = types.StringValue(auth.Kind)
		value["hash"] = types.StringValue(auth.Hash)

	case "PublicKey":
		auth, err := credential.AsUserAuthCredentialUserPublicKeyCredential()

		if err != nil {
			return types.ObjectNull(credentialsAttributes), err
		}

		value["kind"] = types.StringValue(auth.Kind)
		value["public_key"] = types.StringValue(auth.Key)

	case "Sso":
		auth, err := credential.AsUserAuthCredentialUserSsoCredential()

		if err != nil {
			return types.ObjectNull(credentialsAttributes), err
		}

		value["kind"] = types.StringValue(auth.Kind)
		value["email"] = types.StringValue(auth.Email)
		if auth.Provider != nil {
			value["provider"] = types.StringValue(*auth.Provider)
		}

	case "Totp":
		auth, err := credential.AsUserAuthCredentialUserTotpCredential()

		if err != nil {
			return types.ObjectNull(credentialsAttributes), err
		}

		value["kind"] = types.StringValue(auth.Kind)
		value["totp_key"] = ArrayOfUint16ToTerraformList(auth.Key)

	default:
		return types.ObjectNull(credentialsAttributes), errors.New("unknown discriminator value: " + discriminator)
	}

	result = types.ObjectValueMust(credentialsAttributes, value)
	return
}

func ParseUser(user *warpgate.User) (result *provider_models.User, err error) {

	result = &provider_models.User{
		Id:       types.StringValue(user.Id.String()),
		Username: types.StringValue(user.Username),
		Roles:    ArrayOfStringToTerraformSet(user.Roles),
	}

	if len(user.Credentials) == 0 {
		result.Credentials = types.SetNull(types.ObjectType{AttrTypes: credentialsAttributes})
	} else {
		var userCredentials []attr.Value //[]models.UserAuthCredential

		for _, c := range user.Credentials {
			credential, err := ParseUserCredential(c)

			if err != nil {
				return nil, err // todo maybe continue and not return
			}

			userCredentials = append(userCredentials, credential)
		}

		result.Credentials = types.SetValueMust(types.ObjectType{AttrTypes: credentialsAttributes}, userCredentials)
	}

	return
}

func GenerateWarpgateUserAuthCredentials(ctx context.Context, user provider_models.User) (result []warpgate.UserAuthCredential) {

	credentials, err := user.CredentialsAsArray(ctx)

	if err != nil {
		return
	}

	for _, c := range credentials {

		var credential warpgate.UserAuthCredential

		var provider *string

		if c.Provider.IsNull() {
			provider = nil
		} else {
			extractedProvider := c.Provider.ValueString()
			provider = &extractedProvider
		}

		if c.Kind.ValueString() == string(warpgate.Sso) {
			credential.FromUserAuthCredentialUserSsoCredential(
				warpgate.UserAuthCredentialUserSsoCredential{
					Kind:     c.Kind.ValueString(),
					Email:    c.Email.ValueString(),
					Provider: provider, // TODO check for null
				},
			)
		} else if c.Kind.ValueString() == string(warpgate.Password) {
			credential.FromUserAuthCredentialUserPasswordCredential(
				warpgate.UserAuthCredentialUserPasswordCredential{
					Kind: c.Kind.ValueString(),
					Hash: c.Hash.ValueString(),
				},
			)
		} else if c.Kind.ValueString() == string(warpgate.PublicKey) {
			credential.FromUserAuthCredentialUserPublicKeyCredential(
				warpgate.UserAuthCredentialUserPublicKeyCredential{
					Kind: c.Kind.ValueString(),
					Key:  c.PublicKey.ValueString(),
				},
			)
		} else if c.Kind.ValueString() == string(warpgate.Totp) {
			credential.FromUserAuthCredentialUserTotpCredential(
				warpgate.UserAuthCredentialUserTotpCredential{
					Kind: c.Kind.ValueString(),
					Key:  TerraformListToArrayOfUint16(c.TotpKey),
				},
			)
		}

		result = append(result, credential)
	}

	return
}
