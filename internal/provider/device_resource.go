// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

type deviceResource struct {
	client *landb.Client
}

func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*landb.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *deviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *deviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"serialNumber": schema.StringAttribute{
				Optional: true,
			},
			"inventoryNumber": schema.StringAttribute{
				Optional: true,
			},
			"tag": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"zone": schema.StringAttribute{
				Required: true,
			},
			"dhcpResponse": schema.StringAttribute{
				Required: true,
			},
			"ipv4InDnsAndFirewall": schema.BoolAttribute{
				Required: true,
			},
			"ipv6InDnsAndFirewall": schema.BoolAttribute{
				Required: true,
			},
			"managerLock": schema.StringAttribute{
				Required: true,
			},
			"ownership": schema.StringAttribute{
				Required: true,
			},
			"location": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"building": schema.StringAttribute{
						Required: true,
					},
					"floor": schema.StringAttribute{
						Required: true,
					},
					"room": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"parent": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"manufacturer": schema.StringAttribute{
				Required: true,
			},
			"model": schema.StringAttribute{
				Required: true,
			},
			"operatingSystem": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"family": schema.StringAttribute{
						Required: true,
					},
					"version": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"manager": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"person": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
							"username": schema.StringAttribute{
								Required: true,
							},
							"department": schema.StringAttribute{
								Required: true,
							},
							"group": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"egroup": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"reserved": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
			"responsible": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"person": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
							"username": schema.StringAttribute{
								Required: true,
							},
							"department": schema.StringAttribute{
								Required: true,
							},
							"group": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"egroup": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"reserved": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
			"user": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"person": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
							"username": schema.StringAttribute{
								Required: true,
							},
							"department": schema.StringAttribute{
								Required: true,
							},
							"group": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"egroup": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required: true,
							},
							"email": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"reserved": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"firstName": schema.StringAttribute{
								Required: true,
							},
							"lastName": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
			"version": schema.Int64Attribute{
				Required: true,
			},
			"_createdAt": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was created, stored as RFC3339 string",
			},
			"_updatedAt": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was updates, stored as RFC3339 string",
			},
			"_macAddresses": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of MAC addresses associated with the device",
			},
		},
	}
}

func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var deviceState DeviceResource
	diags := req.State.Get(ctx, &deviceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.GetDevice(deviceState.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Device",
			fmt.Sprintf("Could not read Device with Name %s: %s", deviceState.Name, err.Error()),
		)
		return
	}
	if device.SerialNumber != nil {
		serialNumber := types.StringPointerValue(device.SerialNumber)
		deviceState.SerialNumber = &serialNumber
	} else {
		deviceState.SerialNumber = nil
	}

	if device.InventoryNumber != nil {
		inventoryNumber := types.StringPointerValue(device.InventoryNumber)
		deviceState.InventoryNumber = &inventoryNumber
	} else {
		deviceState.InventoryNumber = nil
	}

	if device.Description != nil {
		description := types.StringPointerValue(device.Description)
		deviceState.Description = &description
	} else {
		deviceState.Description = nil
	}

	deviceState.Name = types.StringValue(device.Name)
	deviceState.Tag = types.StringValue(device.Tag)
	deviceState.Zone = types.StringValue(device.Zone)
	deviceState.DHCPResponse = types.StringValue(device.DHCPResponse)
	deviceState.IPv4InDnsFirewall = types.BoolValue(device.IPv4InDnsFirewall)
	deviceState.IPv6InDnsFirewall = types.BoolValue(device.IPv6InDnsFirewall)
	deviceState.ManagerLock = types.StringValue(device.ManagerLock)
	deviceState.Ownership = types.StringValue(device.Ownership)
	deviceState.Location = Location{
		Building: types.StringValue(device.Location.Building),
		Floor:    types.StringValue(device.Location.Floor),
		Room:     types.StringValue(device.Location.Room),
	}
	deviceState.Parent = types.StringValue(device.Parent)
	deviceState.Type = types.StringValue(device.Type)
	deviceState.Manufacturer = types.StringValue(device.Manufacturer)
	deviceState.Model = types.StringValue(device.Model)
	deviceState.OperatingSystem = OperatingSystem{
		Family:  types.StringValue(device.OperatingSystem.Family),
		Version: types.StringValue(device.OperatingSystem.Version),
	}
	deviceState.Manager = Contact{
		Type:     types.StringValue(device.Manager.Type),
		Person:   convertPerson(device.Manager.Person),
		EGroup:   convertEGroup(device.Manager.EGroup),
		Reserved: convertReserved(device.Manager.Reserved),
	}
	deviceState.Responsible = Contact{
		Type:     types.StringValue(device.Responsible.Type),
		Person:   convertPerson(device.Responsible.Person),
		EGroup:   convertEGroup(device.Responsible.EGroup),
		Reserved: convertReserved(device.Responsible.Reserved),
	}
	deviceState.User = Contact{
		Type:     types.StringValue(device.User.Type),
		Person:   convertPerson(device.User.Person),
		EGroup:   convertEGroup(device.User.EGroup),
		Reserved: convertReserved(device.User.Reserved),
	}
	deviceState.Version = types.Int32Value(int32(device.Version))

	diags = resp.State.Set(ctx, deviceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var devicePlan DeviceResource
	diags := req.Plan.Get(ctx, &devicePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device := Device{
		Name:              devicePlan.Name,
		SerialNumber:      devicePlan.SerialNumber,
		InventoryNumber:   devicePlan.InventoryNumber,
		Tag:               devicePlan.Tag,
		Description:       devicePlan.Description,
		Zone:              devicePlan.Zone,
		DHCPResponse:      devicePlan.DHCPResponse,
		IPv4InDnsFirewall: devicePlan.IPv4InDnsFirewall,
		IPv6InDnsFirewall: devicePlan.IPv6InDnsFirewall,
		ManagerLock:       devicePlan.ManagerLock,
		Ownership:         devicePlan.Ownership,
		Location: Location{
			Building: devicePlan.Location.Building,
			Floor:    devicePlan.Location.Floor,
			Room:     devicePlan.Location.Room,
		},
		Parent:       devicePlan.Parent,
		Type:         devicePlan.Type,
		Manufacturer: devicePlan.Manufacturer,
		Model:        devicePlan.Model,
		OperatingSystem: OperatingSystem{
			Family:  devicePlan.OperatingSystem.Family,
			Version: devicePlan.OperatingSystem.Version,
		},
		Manager: Contact{
			Type:     devicePlan.Manager.Type,
			Person:   devicePlan.Manager.Person,
			EGroup:   devicePlan.Manager.EGroup,
			Reserved: devicePlan.Manager.Reserved,
		},
		Responsible: Contact{
			Type:     devicePlan.Responsible.Type,
			Person:   devicePlan.Responsible.Person,
			EGroup:   devicePlan.Responsible.EGroup,
			Reserved: devicePlan.Responsible.Reserved,
		},
		User: Contact{
			Type:     devicePlan.User.Type,
			Person:   devicePlan.User.Person,
			EGroup:   devicePlan.User.EGroup,
			Reserved: devicePlan.User.Reserved,
		},
		Version: devicePlan.Version,
		// UpdatedAt:    devicePlan.UpdatedAt,
	}

	_, err := r.client.UpdateDevice(device.Name.String(), device)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Device",
			fmt.Sprintf("Could not update Device %s: %s", devicePlan.Name, err.Error()),
		)
		return
	}

	device, err = r.client.GetDevice(device.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Device",
			fmt.Sprintf("Could not read Device %s: %s", devicePlan.Name, err.Error()),
		)
		return
	}

	state := DeviceResource{
		ID:                device.Name,
		// ID:                uudi,
		Name:              device.Name,
		SerialNumber:      device.SerialNumber,
		InventoryNumber:   device.InventoryNumber,
		Tag:               device.Tag,
		Description:       device.Description,
		Zone:              device.Zone,
		DHCPResponse:      device.DHCPResponse,
		IPv4InDnsFirewall: device.IPv4InDnsFirewall,
		IPv6InDnsFirewall: device.IPv6InDnsFirewall,
		ManagerLock:       device.ManagerLock,
		Ownership:         device.Ownership,
		Location: Location{
			Building: device.Location.Building,
			Floor:    device.Location.Floor,
			Room:     device.Location.Room,
		},
		Parent:       device.Parent,
		Type:         device.Type,
		Manufacturer: device.Manufacturer,
		Model:        device.Model,
		OperatingSystem: OperatingSystem{
			Family:  device.OperatingSystem.Family,
			Version: device.OperatingSystem.Version,
		},
		Manager: Contact{
			Type:     device.Manager.Type,
			Person:   device.Manager.Person,
			EGroup:   device.Manager.EGroup,
			Reserved: device.Manager.Reserved,
		},
		Responsible: Contact{
			Type:     device.Responsible.Type,
			Person:   device.Responsible.Person,
			EGroup:   device.Responsible.EGroup,
			Reserved: device.Responsible.Reserved,
		},
		User: Contact{
			Type:     device.User.Type,
			Person:   device.User.Person,
			EGroup:   device.User.EGroup,
			Reserved: device.User.Reserved,
		},
		Version: device.Version,
		// UpdatedAt: types.StringValue(time.Now().Format(time.RFC850)),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var deviceState DeviceResource
	diags := req.State.Get(ctx, &deviceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDevice(deviceState.Name.ValueString(), int(deviceState.Version.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Device",
			fmt.Sprintf("Could not delete Device %s: %s", deviceState.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type Device struct {
	Name              types.String    `tfsdk:"name"`
	SerialNumber      *types.String   `tfsdk:"serialNumber"`
	InventoryNumber   *types.String   `tfsdk:"inventoryNumber"`
	Tag               types.String    `tfsdk:"tag"`
	Description       *types.String   `tfsdk:"description"`
	Zone              types.String    `tfsdk:"zone"`
	DHCPResponse      types.String    `tfsdk:"dhcpResponse"`
	IPv4InDnsFirewall types.Bool      `tfsdk:"ipv4InDnsAndFirewall"`
	IPv6InDnsFirewall types.Bool      `tfsdk:"ipv6InDnsAndFirewall"`
	ManagerLock       types.String    `tfsdk:"managerLock"`
	Ownership         types.String    `tfsdk:"ownership"`
	Location          Location        `tfsdk:"location"`
	Parent            types.String    `tfsdk:"parent"`
	Type              types.String    `tfsdk:"type"`
	Manufacturer      types.String    `tfsdk:"manufacturer"`
	Model             types.String    `tfsdk:"model"`
	OperatingSystem   OperatingSystem `tfsdk:"operatingSystem"`
	Manager           Contact         `tfsdk:"manager"`
	Responsible       Contact         `tfsdk:"responsible"`
	User              Contact         `tfsdk:"user"`
	Version           types.Int32     `tfsdk:"version"`
}

type DeviceResource struct {
	ID                types.String    `tfsdk:"id"`
	Name              types.String    `tfsdk:"name"`
	SerialNumber      *types.String   `tfsdk:"serialNumber"`
	InventoryNumber   *types.String   `tfsdk:"inventoryNumber"`
	Tag               types.String    `tfsdk:"tag"`
	Description       *types.String   `tfsdk:"description"`
	Zone              types.String    `tfsdk:"zone"`
	DHCPResponse      types.String    `tfsdk:"dhcpResponse"`
	IPv4InDnsFirewall types.Bool      `tfsdk:"ipv4InDnsAndFirewall"`
	IPv6InDnsFirewall types.Bool      `tfsdk:"ipv6InDnsAndFirewall"`
	ManagerLock       types.String    `tfsdk:"managerLock"`
	Ownership         types.String    `tfsdk:"ownership"`
	Location          Location        `tfsdk:"location"`
	Parent            types.String    `tfsdk:"parent"`
	Type              types.String    `tfsdk:"type"`
	Manufacturer      types.String    `tfsdk:"manufacturer"`
	Model             types.String    `tfsdk:"model"`
	OperatingSystem   OperatingSystem `tfsdk:"operatingSystem"`
	Manager           Contact         `tfsdk:"manager"`
	Responsible       Contact         `tfsdk:"responsible"`
	User              Contact         `tfsdk:"user"`
	Version           types.Int32     `tfsdk:"version"`
	CreatedAt         time.Time       `tfsdk:"_createdAt"`
	UpdatedAt         time.Time       `tfsdk:"_updatedAt"`
	MacAddresses      []types.String  `tfsdk:"_macAddresses"`
}

type DeviceCreateOptions struct {
}

type Location struct {
	Building types.String `tfsdk:"building"`
	Floor    types.String `tfsdk:"floor"`
	Room     types.String `tfsdk:"room"`
}

type OperatingSystem struct {
	Family  types.String `tfsdk:"family"`
	Version types.String `tfsdk:"version"`
}

type Person struct {
	FirstName  types.String `tfsdk:"firstName"`
	LastName   types.String `tfsdk:"lastName"`
	Email      types.String `tfsdk:"email"`
	Username   types.String `tfsdk:"username"`
	Department types.String `tfsdk:"department"`
	Group      types.String `tfsdk:"group"`
}

type Contact struct {
	Type     types.String `tfsdk:"type"`
	Person   *Person      `tfsdk:"person"`
	EGroup   *EGroup      `tfsdk:"egroup"`
	Reserved *Reserved    `tfsdk:"reserved"`
}

type EGroup struct {
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

type Reserved struct {
	FirstName types.String `tfsdk:"firstName"`
	LastName  types.String `tfsdk:"lastName"`
}
