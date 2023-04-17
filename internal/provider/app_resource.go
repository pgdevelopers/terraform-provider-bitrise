package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

// Response from calling register endpoint
type RegisterResponse struct {
	Status  string `json:"status"`
	AppSlug string `json:"slug"`
}

type FinishResponse struct {
	Status            string `json:"status"`
	BuildTriggerToken string `json:"build_trigger_token"`
	BranchName        string `json:"branch_name"`
	Webhook           bool   `json:"is_webhook_auto_reg_supported"`
	WorkflowID        string `json:"default_workflow_id"`
}

type Register struct {
	RepoProvider     string `json:"provider"`
	IsPublic         bool   `json:"is_public"`
	OrganizationSlug string `json:"organization_slug"`
	RepoUrl          string `json:"repo_url"`
	Type             string `json:"type"`
	GitRepoSlug      string `json:"git_repo_slug"`
	GitOwner         string `json:"git_owner"`
	Title            string `json:"title"`
}

type Finish struct {
	ProjectType      string `json:"project_type"`
	StackID          string `json:"stack_id"`
	Config           string `json:"config"`
	Mode             string `json:"mode"`
	OrganizationSlug string `json:"organization_slug"`
}

// AppResourceModel describes the resource data model.
type AppResourceModel struct {
	RepoProvider     types.String `tfsdk:"repo_provider"`
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
			"repo_provider": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Repo provider",
				Computed:            true,
				Default:             stringdefault.StaticString("github"),
			},
			"is_public": schema.BoolAttribute{
				MarkdownDescription: "Is the app public or private",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"organization_slug": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "SLUG for the organization",
				Default:             stringdefault.StaticString("cf38e3d194d03fa2"),
			},
			"repo_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL for the git repository",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Type of the repository",
				Default:             stringdefault.StaticString("git"),
			},
			"git_repo_slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the git repository",
			},
			"git_owner": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
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
			"mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
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

	//***************************** API CALL *******************************
	slug, err := register(data)
	if err != nil {
		return
	}
	_, err = finish(data, slug)
	//*****************************

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

func register(a *AppResourceModel) (string, error) {
	respStruct := RegisterResponse{}
	register := Register{
		RepoProvider:     a.RepoProvider.ValueString(),
		IsPublic:         a.IsPublic.ValueBool(),
		OrganizationSlug: a.OrganizationSlug.ValueString(),
		RepoUrl:          a.RepoUrl.ValueString(),
		Type:             a.Type.ValueString(),
		GitRepoSlug:      a.GitRepoSlug.ValueString(),
		GitOwner:         a.GitOwner.ValueString(),
		Title:            a.Title.ValueString(),
	}
	//resp.Diagnostics.AddError("Look here!", fmt.Sprintf("This is something: %s", register.OrganizationSlug))
	marshalled, err := json.Marshal(register)
	if err != nil {
		return "error", err
	}
	request, err := http.NewRequest("POST", "https://api.bitrise.io/v0.1/apps/register", bytes.NewReader(marshalled))
	if err != nil {
		return "error", err
	}
	request.Header.Set("Authorization", "EZgewzA9KET4uj4cFqoadeLiHwBMKV4orgmZ7kd3AGy_yiMKGBPt050u7KT7fFRd7otH3KGuDKBeftVj0pCxkw")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		return "error", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "error", err
	}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "error", err
	}
	defer res.Body.Close()
	return respStruct.AppSlug, nil
}

func finish(a *AppResourceModel, slug string) (FinishResponse, error) {
	respStruct := FinishResponse{}
	finish := Finish{
		ProjectType:      a.ProjectType.ValueString(),
		Config:           a.Config.ValueString(),
		StackID:          a.StackID.ValueString(),
		OrganizationSlug: a.OrganizationSlug.ValueString(),
		Mode:             a.Mode.ValueString(),
	}
	marshalled, err := json.Marshal(finish)
	if err != nil {
		return respStruct, err
	}
	url := "https://api.bitrise.io/v0.1/apps/" + slug + "/finish"
	request, err := http.NewRequest("POST", url, bytes.NewReader(marshalled))
	if err != nil {
		return respStruct, err
	}
	request.Header.Set("Authorization", "EZgewzA9KET4uj4cFqoadeLiHwBMKV4orgmZ7kd3AGy_yiMKGBPt050u7KT7fFRd7otH3KGuDKBeftVj0pCxkw")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		return respStruct, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return respStruct, err
	}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return respStruct, err
	}
	defer res.Body.Close()
	return respStruct, nil
}
