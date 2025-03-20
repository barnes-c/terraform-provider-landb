package landb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &deviceInterfaceResource{}
	_ resource.ResourceWithConfigure   = &deviceInterfaceResource{}
	_ resource.ResourceWithImportState = &deviceInterfaceResource{}
)

func NewDeviceInterfaceResource() resource.Resource {
	return &deviceInterfaceResource{}
}

type deviceInterfaceResource struct {
	client *landb.Client
}

func (r *deviceInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deviceInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deviceInterface"
}

func (r *deviceInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_name": schema.StringAttribute{
				Required:    true,
				Description: "Virtual machine host name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"interface_domain": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("cern.ch"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_cluster_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_interface_options": schema.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *deviceInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceInterfaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceInterfaceModel := DeviceInterfaceModel{
		DeviceName:             plan.DeviceName,
		InterfaceDomain:        plan.InterfaceDomain,
		DeviceClusterName:      plan.DeviceClusterName,
		DeviceInterfaceOptions: plan.DeviceInterfaceOptions,
	}

	done, err := r.client.DeviceCreate(ctx, deviceInterfaceModel)
	if err != nil || !done {
		resp.Diagnostics.AddError(
			"Error creating Device",
			fmt.Sprintf("Error creating Device %s: %s", deviceInterfaceModel.DeviceName, err.Error()),
		)
		return
	}

	plan.ID = deviceInterfaceModel.DeviceName
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceInterfaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.GetDeviceInterface(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Device Interface",
			fmt.Sprintf("Could not read Device Interface with ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.DeviceName = device.DeviceName
	state.InterfaceDomain = device.InterfaceDomain
	state.DeviceClusterName = device.DeviceClusterName
	state.DeviceInterfaceOptions = device.DeviceInterfaceOptions
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceInterfaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceInterfaceModel := DeviceInterfaceModel{
		DeviceName:             plan.DeviceName,
		InterfaceDomain:        plan.InterfaceDomain,
		DeviceClusterName:      plan.DeviceClusterName,
		DeviceInterfaceOptions: plan.DeviceInterfaceOptions,
	}

	_, err := r.client.UpdateDeviceInterface(ctx, plan.ID.ValueString(), deviceInterfaceModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Device Interface",
			fmt.Sprintf("Could not update Device Interface %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	device, err := r.client.GetDeviceInterface(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Device Interface",
			fmt.Sprintf("Could not read Device Interface %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	state := DeviceInterfaceResourceModel{
		ID:                     device.DeviceName,
		DeviceName:             device.DeviceName,
		InterfaceDomain:        device.InterfaceDomain,
		DeviceClusterName:      device.DeviceClusterName,
		DeviceInterfaceOptions: device.DeviceInterfaceOptions,
		LastUpdated:            types.StringValue(time.Now().Format(time.RFC850)),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceInterfaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDeviceInterface(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Device Interface",
			fmt.Sprintf("Could not delete Device Interface %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *deviceInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type DeviceInterfaceModel struct {
	DeviceName             types.String `tfsdk:"device_name"`
	InterfaceDomain        types.String `tfsdk:"interface_domain"`
	DeviceClusterName      types.String `tfsdk:"defice_cluster_name"`
	DeviceInterfaceOptions types.String `tfsdk:"device_interface_options"`
}

type DeviceInterfaceResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	DeviceName             types.String `tfsdk:"device_name"`
	InterfaceDomain        types.String `tfsdk:"interface_domain"`
	DeviceClusterName      types.String `tfsdk:"defice_cluster_name"`
	DeviceInterfaceOptions types.String `tfsdk:"device_interface_options"`
	LastUpdated            types.String `tfsdk:"last_updated"`
}
