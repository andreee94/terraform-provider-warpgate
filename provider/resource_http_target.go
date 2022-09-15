package provider

import (
	"context"
	"fmt"
	"regexp"
	provider_models "terraform-provider-warpgate/provider/models"
	"terraform-provider-warpgate/provider/validators"
	"terraform-provider-warpgate/warpgate"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces

// var _ provider.ResourceType = httpTargetResourceType{}
var _ resource.Resource = &httpTargetResource{}
var _ resource.ResourceWithImportState = &httpTargetResource{}

// type httpTargetResourceType struct{}

func (r *httpTargetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// "allow_roles": {
			// 	// Type: schema.TypeList,
			// 	// Elem: &schema.Schema{
			// 	// 	Type: schema.TypeString,
			// 	// },
			// 	Type:     types.ListType{ElemType: types.StringType},
			// 	Computed: true,
			// 	Required: false,
			// },
			"id": {
				Computed:            true,
				MarkdownDescription: "Id of the http target in warpgate",
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
				Optional: false,
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"external_host": {
						Type:     types.StringType,
						Computed: false,
						Required: false,
						Optional: true,
						Validators: []tfsdk.AttributeValidator{
							validators.StringRegex{Regex: regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3})$|^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`)},
						},
					},
					"url": {
						Type:     types.StringType,
						Computed: false,
						Required: true,
						Optional: false,
						Validators: []tfsdk.AttributeValidator{
							validators.StringRegex{Regex: regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3})$|^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`)},
						},
					},
					"headers": {
						Type:     types.MapType{ElemType: types.StringType},
						Computed: false,
						Optional: true,
						Required: false,
					},
					"tls": {
						Computed: false,
						Optional: false,
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"mode": {
								Type:     types.StringType,
								Computed: false,
								Required: true,
								Optional: false,
								Validators: []tfsdk.AttributeValidator{
									validators.StringIn([]string{string(warpgate.Disabled), string(warpgate.Preferred), string(warpgate.Required)}, false),
								},
							},
							"verify": {
								Type:     types.BoolType,
								Computed: false,
								Required: true,
								Optional: false,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

// func (t httpTargetResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
// 	provider, diags := convertProviderType(in)

//		return httpTargetResource{
//			provider: provider,
//		}, diags
//	}

func NewHttpTargetResource() resource.Resource {
	return &httpTargetResource{}
}

func (r *httpTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_target"
}

func (r *httpTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type httpTargetResource struct {
	provider *warpgateProvider
}

func (r *httpTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourceState provider_models.TargetHttp

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.provider.client.CreateTargetWithResponse(ctx, warpgate.CreateTargetJSONBody{
		Name: resourceState.Name,
		Options: warpgate.TargetOptionsTargetHTTPOptions{
			Kind:         "Http",
			ExternalHost: &resourceState.Options.ExternalHost.Value,
			Url:          resourceState.Options.Url,
			Headers:      (*warpgate.TargetOptionsTargetHTTPOptions_Headers)(resourceState.Options.Headers),
			Tls: warpgate.Tls{
				Mode:   warpgate.TlsMode(resourceState.Options.Tls.Mode),
				Verify: resourceState.Options.Tls.Verify,
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create http target",
			"Failed to create http target",
		)
		return
	}

	if response.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed to create http target, wrong error code.",
			fmt.Sprintf("Failed to create http target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	resourceState.Id = types.String{Value: response.JSON201.Id.String()}
	// resourceState.AllowRoles = response.JSON201.AllowRoles

	// TODO maybe do not save the password into the state

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *httpTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourceState provider_models.TargetHttp

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
			"Failed to read http target",
			fmt.Sprintf("Failed to read http target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() == 404 {
		resp.Diagnostics.AddWarning(
			"Failed to read http target, resource not found. Removing from the state.",
			fmt.Sprintf("Failed to read http target. (Error code: %d)", response.StatusCode()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to read http target, wrong error code.",
			fmt.Sprintf("Failed to read http target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	httpoptions, err := ParseHttpOptions(response.JSON200.Options)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read http target. Wrong options",
			fmt.Sprintf("Failed to read http target %v. Wrong options type. (Error: %v ", response.JSON200, err),
		)
		return
	}

	// resourceState.AllowRoles = response.JSON200.AllowRoles
	resourceState.Name = response.JSON200.Name
	resourceState.Options.ExternalHost = httpoptions.ExternalHost
	resourceState.Options.Headers = httpoptions.Headers
	resourceState.Options.Tls = httpoptions.Tls
	resourceState.Options.Url = httpoptions.Url

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
}

func (r *httpTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourcePlan provider_models.TargetHttp

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

	var headers *warpgate.TargetOptionsTargetHTTPOptions_Headers

	if resourcePlan.Options.Headers == nil {
		headers = nil
	} else {
		headers = &warpgate.TargetOptionsTargetHTTPOptions_Headers{
			AdditionalProperties: resourcePlan.Options.Headers.AdditionalProperties,
		}
	}

	response, err := r.provider.client.UpdateTargetWithResponse(ctx, id_as_uuid, warpgate.UpdateTargetJSONBody{
		Name: resourcePlan.Name,
		Options: warpgate.TargetOptionsTargetHTTPOptions{
			Kind:         "Http",
			ExternalHost: &resourcePlan.Options.ExternalHost.Value,
			Url:          resourcePlan.Options.Url,
			Headers:      headers,
			Tls: warpgate.Tls{
				Mode:   warpgate.TlsMode(resourcePlan.Options.Tls.Mode),
				Verify: resourcePlan.Options.Tls.Verify,
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update http target",
			fmt.Sprintf("Failed to update http target with id '%s'. (Error: %s)", resourcePlan.Id, err),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed to update http target, wrong error code.",
			fmt.Sprintf("Failed to update http target. (Error code: %d)", response.StatusCode()),
		)
		return
	}

	// probably unnecessary check
	if response.JSON200.Id != id_as_uuid || response.JSON200.Name != resourcePlan.Name {
		resp.Diagnostics.AddWarning(
			"Created resource is different from requested.",
			fmt.Sprintf("Created resource is different from requested. Requested: (%s, %s), Created: (%s, %s)",
				response.JSON200.Id, response.JSON200.Name,
				resourcePlan.Id, resourcePlan.Name,
			),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating http_target state: %v", resourcePlan))

	diags = resp.State.Set(ctx, &resourcePlan)
	resp.Diagnostics.Append(diags...)
}

func (r *httpTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourceState provider_models.TargetHttp

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
			"Failed to delete http target",
			fmt.Sprintf("Failed to delete http target with id '%s'. (Error: %s)", resourceState.Id, err),
		)
		return
	}

	if response.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Failed to delete http target, wrong error code.",
			fmt.Sprintf("Failed to delete http target. (Error code: %d)", response.StatusCode()),
		)
		return
	}
}

func (r *httpTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func ParseHttpOptions(options warpgate.TargetOptions) (*provider_models.TargetHttpOptions, error) {
	var result provider_models.TargetHttpOptions
	var httpoptions warpgate.TargetOptionsTargetHTTPOptions
	err := mapstructure.Decode(options, &httpoptions)

	if err != nil {
		return nil, err
	}

	if httpoptions.ExternalHost == nil {
		result.ExternalHost = types.String{Null: true}
	} else {
		result.ExternalHost = types.String{Value: *httpoptions.ExternalHost}
	}

	if httpoptions.Headers == nil {
		result.Headers = nil
	} else {
		result.Headers = &provider_models.TargetHttpOptions_Headers{
			AdditionalProperties: httpoptions.Headers.AdditionalProperties,
		}
	}

	result.Url = httpoptions.Url
	result.Tls = provider_models.TargetTls{
		Mode:   string(httpoptions.Tls.Mode),
		Verify: httpoptions.Tls.Verify,
	}

	return &result, err
}
