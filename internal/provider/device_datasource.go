// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package provider

import (
	"context"
	"time"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type deviceDataSource struct {
	client *landb.Client
}

func NewDeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*landb.Client); ok {
		d.client = client
	}
}

func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lookup an existing device by its name",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name (unique ID) of the device in LANDB",
			},
			"id": schema.StringAttribute{Computed: true},
			"description":              schema.StringAttribute{Computed: true},
			"dhcp_response":            schema.StringAttribute{Computed: true},
			"inventory_number":         schema.StringAttribute{Computed: true},
			"ipv4_in_dns_and_firewall": schema.BoolAttribute{Computed: true},
			"ipv6_in_dns_and_firewall": schema.BoolAttribute{Computed: true},
			"last_updated":             schema.StringAttribute{Computed: true},
			"location": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Physical location of the device",
				Attributes: map[string]schema.Attribute{
					"building": schema.StringAttribute{Computed: true},
					"floor":    schema.StringAttribute{Computed: true},
					"room":     schema.StringAttribute{Computed: true},
				},
			},
			"manager_lock": schema.StringAttribute{Computed: true},
			"manager":     contactSchemaBlock("Manager of the device"),
			"manufacturer": schema.StringAttribute{Computed: true},
			"model":        schema.StringAttribute{Computed: true},
			"ownership":    schema.StringAttribute{Computed: true},
			"parent":       schema.StringAttribute{Computed: true},
			"serial_number": schema.StringAttribute{Computed: true},
			"tag":           schema.StringAttribute{Computed: true},
			"type":          schema.StringAttribute{Computed: true},
			"zone":          schema.StringAttribute{Computed: true},
			"operating_system": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Operating system of the device",
				Attributes: map[string]schema.Attribute{
					"family":  schema.StringAttribute{Computed: true},
					"version": schema.StringAttribute{Computed: true},
				},
			},
			"responsible": contactSchemaBlock("Responsible of the device"),
			"user":        contactSchemaBlock("User of the device"),
			"version": schema.Int64Attribute{Computed: true},
		},
	}
}

func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data deviceResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := d.client.GetDevice(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching device", err.Error())
		return
	}

	data.ID = types.StringValue(device.Name)
	data.Description = types.StringValue(device.Description)
	data.DHCPResponse = types.StringValue(device.DHCPResponse)
	data.InventoryNumber = types.StringValue(device.InventoryNumber)
	data.IPv4InDNSAndFirewall = types.BoolValue(device.IPv4InDNSAndFirewall)
	data.IPv6InDNSAndFirewall = types.BoolValue(device.IPv6InDNSAndFirewall)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	data.Location = flattenLocation(device.Location)
	data.ManagerLock = types.StringValue(device.ManagerLock)
	data.Manager = flattenContactObject(device.Manager)
	data.Manufacturer = types.StringValue(device.Manufacturer)
	data.Model = types.StringValue(device.Model)
	data.Ownership = types.StringValue(device.Ownership)
	data.Parent = types.StringValue(device.Parent)
	data.SerialNumber = types.StringValue(device.SerialNumber)
	data.Tag = types.StringValue(device.Tag)
	data.Type = types.StringValue(device.Type)
	data.Zone = types.StringValue(device.Zone)
	data.OperatingSystem = flattenOperatingSystem(device.OperatingSystem)
	data.Responsible = flattenContactObject(device.Responsible)
	data.User = flattenContactObject(device.User)
	data.Version = types.Int64Value(int64(device.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
