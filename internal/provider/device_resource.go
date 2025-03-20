package landb

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
			fmt.Sprintf("Expected *landb.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_name": schema.StringAttribute{
				Required: true,
			},
			"manufacturer": schema.StringAttribute{
				Required: true,
			},
			"model": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"tag": schema.StringAttribute{
				Optional: true,
			},
			"ipv6_ready": schema.BoolAttribute{
				Required: true,
			},
			"location": schema.SingleNestedAttribute{
				Required: true,
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
			"operating_system": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"version": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"landb_manager_person": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"first_name": schema.StringAttribute{
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
			"responsible_person": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"first_name": schema.StringAttribute{
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
			"user_person": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"first_name": schema.StringAttribute{
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
		},
	}
}

func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceModel := DeviceModel{
		DeviceName:   plan.DeviceName,
		Manufacturer: plan.Manufacturer,
		Model:        plan.Model,
		Description:  plan.Description,
		Tag:          plan.Tag,
		IPv6Ready:    plan.IPv6Ready,
		Location: LocationModel{
			Building: plan.Location.Building,
			Floor:    plan.Location.Floor,
			Room:     plan.Location.Room,
		},
		OperatingSystem: OperatingSystemModel{
			Name:    plan.OperatingSystem.Name,
			Version: plan.OperatingSystem.Version,
		},
		LandbManagerPerson: PersonModel{
			Name:       plan.LandbManagerPerson.Name,
			FirstName:  plan.LandbManagerPerson.FirstName,
			Department: plan.LandbManagerPerson.Department,
			Group:      plan.LandbManagerPerson.Group,
		},
		ResponsiblePerson: PersonModel{
			Name:       plan.ResponsiblePerson.Name,
			FirstName:  plan.ResponsiblePerson.FirstName,
			Department: plan.ResponsiblePerson.Department,
			Group:      plan.ResponsiblePerson.Group,
		},
		UserPerson: PersonModel{
			Name:       plan.UserPerson.Name,
			FirstName:  plan.UserPerson.FirstName,
			Department: plan.UserPerson.Department,
			Group:      plan.UserPerson.Group,
		},
	}

	createOptions := DeviceCreateOptions{}

	done, err := r.client.DeviceCreate(ctx, deviceModel, createOptions)
	if err != nil || !done {
		resp.Diagnostics.AddError(
			"Error creating Device",
			fmt.Sprintf("Error creating Device %s: %s", deviceModel.DeviceName, err.Error()),
		)
		return
	}

	plan.ID = deviceModel.DeviceName
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read fetches the latest state from the API and updates Terraform.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the Device from the API.
	device, err := r.client.GetDevice(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Device",
			fmt.Sprintf("Could not read Device with ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.DeviceName = device.DeviceName
	state.Manufacturer = device.Manufacturer
	state.Model = device.Model
	state.Description = device.Description
	state.Tag = device.Tag
	state.IPv6Ready = device.IPv6Ready
	state.Location = LocationModel{
		Building: device.Location.Building,
		Floor:    device.Location.Floor,
		Room:     device.Location.Room,
	}
	state.OperatingSystem = OperatingSystemModel{
		Name:    device.OperatingSystem.Name,
		Version: device.OperatingSystem.Version,
	}
	state.LandbManagerPerson = PersonModel{
		Name:       device.LandbManagerPerson.Name,
		FirstName:  device.LandbManagerPerson.FirstName,
		Department: device.LandbManagerPerson.Department,
		Group:      device.LandbManagerPerson.Group,
	}
	state.ResponsiblePerson = PersonModel{
		Name:       device.ResponsiblePerson.Name,
		FirstName:  device.ResponsiblePerson.FirstName,
		Department: device.ResponsiblePerson.Department,
		Group:      device.ResponsiblePerson.Group,
	}
	state.UserPerson = PersonModel{
		Name:       device.UserPerson.Name,
		FirstName:  device.UserPerson.FirstName,
		Department: device.UserPerson.Department,
		Group:      device.UserPerson.Group,
	}
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceModel := DeviceModel{
		DeviceName:   plan.DeviceName,
		Manufacturer: plan.Manufacturer,
		Model:        plan.Model,
		Description:  plan.Description,
		Tag:          plan.Tag,
		IPv6Ready:    plan.IPv6Ready,
		Location: LocationModel{
			Building: plan.Location.Building,
			Floor:    plan.Location.Floor,
			Room:     plan.Location.Room,
		},
		OperatingSystem: OperatingSystemModel{
			Name:    plan.OperatingSystem.Name,
			Version: plan.OperatingSystem.Version,
		},
		LandbManagerPerson: PersonModel{
			Name:       plan.LandbManagerPerson.Name,
			FirstName:  plan.LandbManagerPerson.FirstName,
			Department: plan.LandbManagerPerson.Department,
			Group:      plan.LandbManagerPerson.Group,
		},
		ResponsiblePerson: PersonModel{
			Name:       plan.ResponsiblePerson.Name,
			FirstName:  plan.ResponsiblePerson.FirstName,
			Department: plan.ResponsiblePerson.Department,
			Group:      plan.ResponsiblePerson.Group,
		},
		UserPerson: PersonModel{
			Name:       plan.UserPerson.Name,
			FirstName:  plan.UserPerson.FirstName,
			Department: plan.UserPerson.Department,
			Group:      plan.UserPerson.Group,
		},
	}

	_, err := r.client.UpdateDevice(ctx, plan.ID.ValueString(), deviceModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Device",
			fmt.Sprintf("Could not update Device %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	device, err := r.client.GetDevice(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Device",
			fmt.Sprintf("Could not read Device %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	state := DeviceResourceModel{
		ID:           device.DeviceName,
		DeviceName:   device.DeviceName,
		Manufacturer: device.Manufacturer,
		Model:        device.Model,
		Description:  device.Description,
		Tag:          device.Tag,
		IPv6Ready:    device.IPv6Ready,
		Location: LocationModel{
			Building: device.Location.Building,
			Floor:    device.Location.Floor,
			Room:     device.Location.Room,
		},
		OperatingSystem: OperatingSystemModel{
			Name:    device.OperatingSystem.Name,
			Version: device.OperatingSystem.Version,
		},
		LandbManagerPerson: PersonModel{
			Name:       device.LandbManagerPerson.Name,
			FirstName:  device.LandbManagerPerson.FirstName,
			Department: device.LandbManagerPerson.Department,
			Group:      device.LandbManagerPerson.Group,
		},
		ResponsiblePerson: PersonModel{
			Name:       device.ResponsiblePerson.Name,
			FirstName:  device.ResponsiblePerson.FirstName,
			Department: device.ResponsiblePerson.Department,
			Group:      device.ResponsiblePerson.Group,
		},
		UserPerson: PersonModel{
			Name:       device.UserPerson.Name,
			FirstName:  device.UserPerson.FirstName,
			Department: device.UserPerson.Department,
			Group:      device.UserPerson.Group,
		},
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDevice(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Device",
			fmt.Sprintf("Could not delete Device %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type DeviceModel struct {
	Description        types.String         `json:"description"`
	DeviceName         types.String         `json:"device_name"`
	IPv6Ready          types.Bool           `json:"ipv6_ready"`
	LandbManagerPerson PersonModel          `json:"landb_manager_person"`
	Location           LocationModel        `json:"location"`
	Manufacturer       types.String         `json:"manufacturer"`
	Model              types.String         `json:"model"`
	OperatingSystem    OperatingSystemModel `json:"operating_system"`
	ResponsiblePerson  PersonModel          `json:"responsible_person"`
	Tag                types.String         `json:"tag"`
	UserPerson         PersonModel          `json:"user_person"`
}

type DeviceResourceModel struct {
	ID                 types.String         `tfsdk:"id"`
	Description        types.String         `tfsdk:"description"`
	DeviceName         types.String         `tfsdk:"device_name"`
	IPv6Ready          types.Bool           `tfsdk:"ipv6_ready"`
	LandbManagerPerson PersonModel          `tfsdk:"landb_manager_person"`
	LastUpdated        types.String         `tfsdk:"last_updated"`
	Location           LocationModel        `tfsdk:"location"`
	Manufacturer       types.String         `tfsdk:"manufacturer"`
	Model              types.String         `tfsdk:"model"`
	OperatingSystem    OperatingSystemModel `tfsdk:"operating_system"`
	ResponsiblePerson  PersonModel          `tfsdk:"responsible_person"`
	Tag                types.String         `tfsdk:"tag"`
	UserPerson         PersonModel          `tfsdk:"user_person"`
}

type DeviceCreateOptions struct {
}

type LocationModel struct {
	Building types.String `tfsdk:"building"`
	Floor    types.String `tfsdk:"floor"`
	Room     types.String `tfsdk:"room"`
}

type OperatingSystemModel struct {
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

type PersonModel struct {
	Name       types.String `tfsdk:"name"`
	FirstName  types.String `tfsdk:"first_name"`
	Department types.String `tfsdk:"department"`
	Group      types.String `tfsdk:"group"`
}
