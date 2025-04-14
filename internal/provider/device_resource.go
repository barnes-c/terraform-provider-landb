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
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

type deviceResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	SerialNumber         types.String `tfsdk:"serial_number"`
	InventoryNumber      types.String `tfsdk:"inventory_number"`
	Tag                  types.String `tfsdk:"tag"`
	Description          types.String `tfsdk:"description"`
	Zone                 types.String `tfsdk:"zone"`
	DHCPResponse         types.String `tfsdk:"dhcp_response"`
	IPv4InDNSAndFirewall types.Bool   `tfsdk:"ipv4_in_dns_and_firewall"`
	IPv6InDNSAndFirewall types.Bool   `tfsdk:"ipv6_in_dns_and_firewall"`
	ManagerLock          types.String `tfsdk:"manager_lock"`
	Ownership            types.String `tfsdk:"ownership"`
	Type                 types.String `tfsdk:"type"`
	Parent               types.String `tfsdk:"parent"`
	Manufacturer         types.String `tfsdk:"manufacturer"`
	Model                types.String `tfsdk:"model"`
	Version              types.Int64  `tfsdk:"version"`
	LastUpdated          types.String `tfsdk:"last_updated"`
}

type deviceResource struct {
	client *landb.Client
}

func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Device name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Device name.",
				Required:    true,
			},
			"serial_number": schema.StringAttribute{
				Description: "Serial number of the device.",
				Optional:    true,
			},
			"inventory_number": schema.StringAttribute{
				Description: "Inventory number of the device.",
				Optional:    true,
			},
			"tag": schema.StringAttribute{
				Description: "Tag of the device.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the device.",
				Optional:    true,
			},
			"zone": schema.StringAttribute{
				Description: "Zone of the device.",
				Required:    true,
			},
			"dhcp_response": schema.StringAttribute{
				Description: "DHCP Response behavior (e.g., ALWAYS).",
				Required:    true,
			},
			"ipv4_in_dns_and_firewall": schema.BoolAttribute{
				Description: "Whether IPv4 is in DNS and Firewall.",
				Required:    true,
			},
			"ipv6_in_dns_and_firewall": schema.BoolAttribute{
				Description: "Whether IPv6 is in DNS and Firewall.",
				Required:    true,
			},
			"manager_lock": schema.StringAttribute{
				Description: "Manager lock status (e.g., NO_LOCK).",
				Required:    true,
			},
			"ownership": schema.StringAttribute{
				Description: "Ownership (e.g., CERN).",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the device (e.g., BRIDGE).",
				Required:    true,
			},
			"parent": schema.StringAttribute{
				Description: "Parent device name.",
				Optional:    true,
			},
			"manufacturer": schema.StringAttribute{
				Description: "Manufacturer of the device.",
				Optional:    true,
			},
			"model": schema.StringAttribute{
				Description: "Model of the device.",
				Optional:    true,
			},
			"version": schema.Int64Attribute{
				Description: "Version for optimistic locking.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of last Terraform update.",
				Computed:    true,
			},
		},
	}
}

func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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
	}

	createdDevice, err := r.client.CreateDevice(device)
	if err != nil {
		resp.Diagnostics.AddError("Error creating device", "Could not create device: "+err.Error())
		return
	}

	plan.ID = types.StringValue(createdDevice.Name)
	plan.Version = types.Int64Value(int64(createdDevice.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.GetDevice(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device", "Could not read device: "+err.Error())
		return
	}

	state.Version = types.Int64Value(int64(device.Version))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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
	}

	updatedSet, err := r.client.UpdateDevice(plan.Name.ValueString(), device)
	if err != nil {
		resp.Diagnostics.AddError("Error updating device", "Could not update device: "+err.Error())
		return
	}

	plan.Version = types.Int64Value(int64(updatedSet.Version))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSet(state.Name.ValueString(), int(state.Version.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting device", "Could not delete device: "+err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
