package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AppResource{}
var _ resource.ResourceWithImportState = &AppResource{}

func NewAppResource() resource.Resource {
	return &AppResource{}
}

// AppResource defines the resource implementation.
type AppResource struct {
	client *http.Client
}

// AppResourceModel describes the resource data model.
type AppResourceModel struct {
	Provider         types.String `tfsdk:"provider"`
	IsPublic         types.Bool   `tfsdk:"is_public"`
	OrganizationSlug types.String `tfsdk:"organization_slug"`
	RepoUrl          types.String `tfsdk:"repo_url"`
	Type             types.String `tfsdk:"type"`
	GitRepoSlug      types.String `tfsdk:"git_repo_slug"`
	GitOwner         types.String `tfsdk:"git_owner"`
	Title            types.String `tfsdk:"title"`
	ProjectType      types.String `tfsdk:"project_type"`
	StackID          types.String `tfsdk:"stack_id"`
	Config           types.String `tfsdk:"config"`
	Mode             types.String `tfsdk:"mode"`
}

func (r *AppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *AppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "App resource",

		Attributes: map[string]schema.Attribute{
			"provider": schema.StringAttribute{
				MarkdownDescription: "Repo provider",
				Optional:            true,
				Default:             stringdefault.StaticString("github"),
			},
			"is_public": schema.BoolAttribute{
				MarkdownDescription: "Is the app public or private",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"organization_slug": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SLUG for the organization",
				Default:             stringdefault.StaticString("cf38e3d194d03fa2"),
			},
			"repo_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL for the git repository",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Type of the repository",
				Default:             stringdefault.StaticString("git"),
			},
			"git_repo_slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the git repository",
			},
			"git_owner": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Github org name",
				Default:             stringdefault.StaticString("pgdevelopers"),
			},
			"title": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional app rename",
			},
			"project_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Operating system",
			},
			"stack_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Not sure?",
			},
			"config": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OS configuration?",
			},
			"Mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Must be manual",
				Default:             stringdefault.StaticString("manual"),
			},
		},
	}
}

func (r *AppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AppResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create App, got error: %s", err))
	//     return
	// }

	// For the purposes of this App code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.StringValue("App-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AppResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read App, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AppResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update App, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AppResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete App, got error: %s", err))
	//     return
	// }
}

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
