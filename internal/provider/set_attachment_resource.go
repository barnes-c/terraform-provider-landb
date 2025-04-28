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

type setAttachmentResourceModel struct {
	ID          types.String `tfsdk:"id"`
	SetName     types.String `tfsdk:"set_name"`
	Name        types.String `tfsdk:"name"`
	IPv4        types.String `tfsdk:"ipv4"`
	IPv6        types.String `tfsdk:"ipv6"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

type setAttachmentResource struct {
	client *landb.Client
}

func NewSetAttachmentResource() resource.Resource {
	return &setAttachmentResource{}
}

func (r *setAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_attach"
}

func (r *setAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a single IP-address attachment on a set",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"set_name": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ipv4": schema.StringAttribute{
				Required: true,
			},
			"ipv6": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *setAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*landb.Client); ok {
		r.client = client
	}
}

func (r *setAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan setAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	att := landb.SetAttachment{
		Name:        plan.Name.ValueString(),
		IPv4:        plan.IPv4.ValueString(),
		IPv6:        plan.IPv6.ValueString(),
		Description: plan.Description.ValueString(),
	}

	created, err := r.client.CreateSetAttachment(plan.SetName.ValueString(), att)
	if err != nil {
		resp.Diagnostics.AddError("Error creating set attachment", err.Error())
		return
	}

	plan.ID = types.StringValue(created.Name)
	plan.CreatedAt = types.StringValue(created.CreatedAt.Format(time.RFC850))
	plan.UpdatedAt = types.StringValue(created.UpdatedAt.Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *setAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state setAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	all, err := r.client.GetSetAttachments(state.SetName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error listing set attachments", err.Error())
		return
	}

	var found *landb.SetAttachment
	for _, a := range all {
		if a.Name == state.ID.ValueString() {
			found = &a
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(found.Name)
	state.IPv4 = types.StringValue(found.IPv4)
	state.IPv6 = types.StringValue(found.IPv6)
	state.Description = types.StringValue(found.Description)
	state.CreatedAt = types.StringValue(found.CreatedAt.Format(time.RFC850))
	state.UpdatedAt = types.StringValue(found.UpdatedAt.Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *setAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan setAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	att := landb.SetAttachment{
		Name:        plan.Name.ValueString(),
		IPv4:        plan.IPv4.ValueString(),
		IPv6:        plan.IPv6.ValueString(),
		Description: plan.Description.ValueString(),
	}

	updated, err := r.client.UpdateSetAttachment(plan.SetName.ValueString(), plan.ID.ValueString(), att)
	if err != nil {
		resp.Diagnostics.AddError("Error updating set attachment", err.Error())
		return
	}

	plan.UpdatedAt = types.StringValue(updated.UpdatedAt.Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *setAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state setAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSetAttachment(state.SetName.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting set attachment", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *setAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
