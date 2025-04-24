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

type deviceResourceModel struct {
	Description          types.String `tfsdk:"description"`
	DHCPResponse         types.String `tfsdk:"dhcp_response"`
	ID                   types.String `tfsdk:"id"`
	InventoryNumber      types.String `tfsdk:"inventory_number"`
	IPv4InDNSAndFirewall types.Bool   `tfsdk:"ipv4_in_dns_and_firewall"`
	IPv6InDNSAndFirewall types.Bool   `tfsdk:"ipv6_in_dns_and_firewall"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Manager              types.Object `tfsdk:"manager"`
	ManagerLock          types.String `tfsdk:"manager_lock"`
	Manufacturer         types.String `tfsdk:"manufacturer"`
	Model                types.String `tfsdk:"model"`
	Name                 types.String `tfsdk:"name"`
	OperatingSystem      types.Object `tfsdk:"operating_system"`
	Ownership            types.String `tfsdk:"ownership"`
	Parent               types.String `tfsdk:"parent"`
	Responsible          types.Object `tfsdk:"responsible"`
	SerialNumber         types.String `tfsdk:"serial_number"`
	Tag                  types.String `tfsdk:"tag"`
	Type                 types.String `tfsdk:"type"`
	User                 types.Object `tfsdk:"user"`
	Version              types.Int64  `tfsdk:"version"`
	Zone                 types.String `tfsdk:"zone"`
}

type deviceResource struct {
	client *landb.Client
}

func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a device",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description":              schema.StringAttribute{Optional: true},
			"dhcp_response":            schema.StringAttribute{Required: true},
			"inventory_number":         schema.StringAttribute{Optional: true},
			"ipv4_in_dns_and_firewall": schema.BoolAttribute{Required: true},
			"ipv6_in_dns_and_firewall": schema.BoolAttribute{Required: true},
			"last_updated":             schema.StringAttribute{Computed: true},
			"manager_lock":             schema.StringAttribute{Required: true},
			"manager":                  contactSchemaBlock("Manager of the device"),
			"manufacturer":             schema.StringAttribute{Optional: true},
			"model":                    schema.StringAttribute{Optional: true},
			"name":                     schema.StringAttribute{Required: true},
			"operating_system": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Operating system of the device",
				Attributes: map[string]schema.Attribute{
					"name":    schema.StringAttribute{Required: true},
					"version": schema.StringAttribute{Optional: true},
				},
			},
			"ownership":     schema.StringAttribute{Required: true},
			"parent":        schema.StringAttribute{Optional: true},
			"responsible":   contactSchemaBlock("Responsible person of the device"),
			"serial_number": schema.StringAttribute{Optional: true},
			"tag":           schema.StringAttribute{Optional: true},
			"type":          schema.StringAttribute{Required: true},
			"user":          contactSchemaBlock("User of the device"),
			"version":       schema.Int64Attribute{Computed: true},
			"zone":          schema.StringAttribute{Required: true},
		},
	}
}

func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*landb.Client); ok {
		r.client = client
	}
}

func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	manager, managerDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(managerDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	responsible, responsibleDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(responsibleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, userDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(userDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	os, osDiags := expandOperatingSystem(ctx, plan.OperatingSystem)
	resp.Diagnostics.Append(osDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device := landb.Device{
		Description:          plan.Description.ValueString(),
		DHCPResponse:         plan.DHCPResponse.ValueString(),
		InventoryNumber:      plan.InventoryNumber.ValueString(),
		IPv4InDNSAndFirewall: plan.IPv4InDNSAndFirewall.ValueBool(),
		IPv6InDNSAndFirewall: plan.IPv6InDNSAndFirewall.ValueBool(),
		Manager:              manager,
		ManagerLock:          plan.ManagerLock.ValueString(),
		Manufacturer:         plan.Manufacturer.ValueString(),
		Model:                plan.Model.ValueString(),
		Name:                 plan.Name.ValueString(),
		OperatingSystem:      os,
		Ownership:            plan.Ownership.ValueString(),
		Parent:               plan.Parent.ValueString(),
		Responsible:          responsible,
		SerialNumber:         plan.SerialNumber.ValueString(),
		Tag:                  plan.Tag.ValueString(),
		Type:                 plan.Type.ValueString(),
		User:                 user,
		Zone:                 plan.Zone.ValueString(),
	}

	created, err := r.client.CreateDevice(device)
	if err != nil {
		resp.Diagnostics.AddError("Error creating device", err.Error())
		return
	}

	plan.ID = types.StringValue(created.Name)
	plan.Version = types.Int64Value(int64(created.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Manager = flattenContactObject(created.Manager)
	plan.Responsible = flattenContactObject(created.Responsible)
	plan.User = flattenContactObject(created.User)
	plan.OperatingSystem = flattenOperatingSystem(created.OperatingSystem)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	devicePtr, err := r.client.GetDevice(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device", err.Error())
		return
	}

	state.Version = types.Int64Value(int64(devicePtr.Version))
	state.Manager = flattenContactObject(devicePtr.Manager)
	state.Responsible = flattenContactObject(devicePtr.Responsible)
	state.User = flattenContactObject(devicePtr.User)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	manager, managerDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(managerDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	responsible, responsibleDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(responsibleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, userDiags := expandContactObject(ctx, plan.Responsible)
	resp.Diagnostics.Append(userDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	os, osDiags := expandOperatingSystem(ctx, plan.OperatingSystem)
	resp.Diagnostics.Append(osDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device := landb.Device{
		Name:                 plan.Name.ValueString(),
		SerialNumber:         plan.SerialNumber.ValueString(),
		InventoryNumber:      plan.InventoryNumber.ValueString(),
		Tag:                  plan.Tag.ValueString(),
		Description:          plan.Description.ValueString(),
		Zone:                 plan.Zone.ValueString(),
		DHCPResponse:         plan.DHCPResponse.ValueString(),
		IPv4InDNSAndFirewall: plan.IPv4InDNSAndFirewall.ValueBool(),
		IPv6InDNSAndFirewall: plan.IPv6InDNSAndFirewall.ValueBool(),
		ManagerLock:          plan.ManagerLock.ValueString(),
		Ownership:            plan.Ownership.ValueString(),
		Type:                 plan.Type.ValueString(),
		Parent:               plan.Parent.ValueString(),
		Manufacturer:         plan.Manufacturer.ValueString(),
		Model:                plan.Model.ValueString(),
		Manager:              manager,
		Responsible:          responsible,
		User:                 user,
		OperatingSystem:      os,
	}

	updated, err := r.client.UpdateDevice(plan.Name.ValueString(), device)
	if err != nil {
		resp.Diagnostics.AddError("Error updating device", err.Error())
		return
	}

	plan.Version = types.Int64Value(int64(updated.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_ = r.client.DeleteDevice(state.Name.ValueString(), int(state.Version.ValueInt64()))
	resp.State.RemoveResource(ctx)
}

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
