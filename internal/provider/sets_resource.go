// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later
package provider

import (
	"context"
	"time"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	NetworkDomain        types.String `tfsdk:"network_domain"`
	Responsible          types.Object `tfsdk:"responsible"`
	Description          types.String `tfsdk:"description"`
	ProjectURL           types.String `tfsdk:"project_url"`
	ReceiveNotifications types.Bool   `tfsdk:"receive_notifications"`
	Version              types.Int64  `tfsdk:"version"`
	LastUpdated          types.String `tfsdk:"last_updated"`
}

type setResource struct {
	client *landb.Client
}

func NewSetResource() resource.Resource {
	return &setResource{}
}

func (r *setResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set"
}

func (r *setResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a set",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":                  schema.StringAttribute{Required: true},
			"type":                  schema.StringAttribute{Required: true},
			"network_domain":        schema.StringAttribute{Required: true},
			"responsible":           contactSchemaBlock("Responsible entity for the set"),
			"description":           schema.StringAttribute{Optional: true},
			"project_url":           schema.StringAttribute{Optional: true},
			"receive_notifications": schema.BoolAttribute{Optional: true},
			"version":               schema.Int64Attribute{Computed: true},
			"last_updated":          schema.StringAttribute{Computed: true},
		},
	}
}

func (r *setResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*landb.Client); ok {
		r.client = client
	}
}

func (r *setResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan setResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	responsible, responsibleDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(responsibleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	setObj := landb.Set{
		Name:                 plan.Name.ValueString(),
		Type:                 plan.Type.ValueString(),
		NetworkDomain:        plan.NetworkDomain.ValueString(),
		Description:          plan.Description.ValueString(),
		ProjectURL:           plan.ProjectURL.ValueString(),
		ReceiveNotifications: plan.ReceiveNotifications.ValueBool(),
		Responsible:          responsible,
	}

	created, err := r.client.CreateSet(setObj)
	if err != nil {
		resp.Diagnostics.AddError("Error creating set", err.Error())
		return
	}

	plan.ID = types.StringValue(created.Name)
	plan.Version = types.Int64Value(int64(created.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Responsible = flattenContactObject(created.Responsible)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *setResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state setResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ptr, err := r.client.GetSet(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading set", err.Error())
		return
	}

	state.Version = types.Int64Value(int64(ptr.Version))
	state.Responsible = flattenContactObject(ptr.Responsible)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *setResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan setResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	responsible, responsibleDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(responsibleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	setObj := landb.Set{
		Name:                 plan.Name.ValueString(),
		Type:                 plan.Type.ValueString(),
		NetworkDomain:        plan.NetworkDomain.ValueString(),
		Description:          plan.Description.ValueString(),
		ProjectURL:           plan.ProjectURL.ValueString(),
		ReceiveNotifications: plan.ReceiveNotifications.ValueBool(),
		Responsible:          responsible,
	}

	updated, err := r.client.UpdateSet(plan.Name.ValueString(), setObj)
	if err != nil {
		resp.Diagnostics.AddError("Error updating set", err.Error())
		return
	}

	plan.Version = types.Int64Value(int64(updated.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *setResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state setResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_ = r.client.DeleteSet(state.Name.ValueString(), int(state.Version.ValueInt64()))
	resp.State.RemoveResource(ctx)
}

func (r *setResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
