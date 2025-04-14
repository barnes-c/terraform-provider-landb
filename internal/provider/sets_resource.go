// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	landb "landb/internal/client"
)

var (
	_ resource.Resource                = &setResource{}
	_ resource.ResourceWithConfigure   = &setResource{}
	_ resource.ResourceWithImportState = &setResource{}
)

func NewSetResource() resource.Resource {
	return &setResource{}
}

type setResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	NetworkDomain        types.String `tfsdk:"network_domain"`
	Description          types.String `tfsdk:"description"`
	ProjectURL           types.String `tfsdk:"project_url"`
	ReceiveNotifications types.Bool   `tfsdk:"receive_notifications"`
	Version              types.Int64  `tfsdk:"version"`
	LastUpdated          types.String `tfsdk:"last_updated"`
}

type setResource struct {
	client *landb.Client
}

func (r *setResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set"
}

func (r *setResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of the set.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the set.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the set (e.g., INTERDOMAIN).",
				Required:    true,
			},
			"network_domain": schema.StringAttribute{
				Description: "Network domain of the set.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the set.",
				Optional:    true,
			},
			"project_url": schema.StringAttribute{
				Description: "Project URL associated with the set.",
				Optional:    true,
			},
			"receive_notifications": schema.BoolAttribute{
				Description: "Whether to receive notifications.",
				Optional:    true,
			},
			"version": schema.Int64Attribute{
				Description: "Version of the set for optimistic locking.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the set.",
				Computed:    true,
			},
		},
	}
}

func (r *setResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan setResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	set := landb.Set{
		Name:                 plan.Name.ValueString(),
		Type:                 plan.Type.ValueString(),
		NetworkDomain:        plan.NetworkDomain.ValueString(),
		Description:          plan.Description.ValueString(),
		ProjectURL:           plan.ProjectURL.ValueString(),
		ReceiveNotifications: plan.ReceiveNotifications.ValueBool(),
	}

	createdSet, err := r.client.CreateSet(set)
	if err != nil {
		resp.Diagnostics.AddError("Error creating set", "Could not create set: "+err.Error())
		return
	}

	plan.ID = types.StringValue(createdSet.Name)
	plan.Version = types.Int64Value(int64(createdSet.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *setResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state setResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetSet(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading set", "Could not read set: "+err.Error())
		return
	}

	state.Version = types.Int64Value(int64(set.Version))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *setResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan setResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	set := landb.Set{
		Name:                 plan.Name.ValueString(),
		Type:                 plan.Type.ValueString(),
		NetworkDomain:        plan.NetworkDomain.ValueString(),
		Description:          plan.Description.ValueString(),
		ProjectURL:           plan.ProjectURL.ValueString(),
		ReceiveNotifications: plan.ReceiveNotifications.ValueBool(),
		Version:              int(plan.Version.ValueInt64()),
	}

	updatedSet, err := r.client.UpdateSet(plan.Name.ValueString(), set)
	if err != nil {
		resp.Diagnostics.AddError("Error updating set", "Could not update set: "+err.Error())
		return
	}

	plan.Version = types.Int64Value(int64(updatedSet.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *setResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state setResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSet(state.Name.ValueString(), int(state.Version.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting set", "Could not delete set: "+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *setResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*landb.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected client type",
			fmt.Sprintf("Expected *landb.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *setResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
