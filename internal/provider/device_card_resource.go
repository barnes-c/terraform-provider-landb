package landb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &deviceCardResource{}
	_ resource.ResourceWithConfigure   = &deviceCardResource{}
	_ resource.ResourceWithImportState = &deviceCardResource{}
)

func NewDeviceCardResource() resource.Resource {
	return &deviceCardResource{}
}

type deviceCardResource struct {
	client *landb.Client
}

func (r *deviceCardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deviceCardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deviceCard"
}

func (r *deviceCardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_name": schema.StringAttribute{
				Required:    true,
				Description: "Device name (VM host name)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hardware_address": schema.StringAttribute{
				Required:    true,
				Description: "Hardware address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"card_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of card (default Ethernet)",
				Default:     stringdefault.StaticString("Ethernet"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *deviceCardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceCardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceCardModel := DeviceCardModel{
		DeviceName:      plan.DeviceName,
		HardwareAddress: plan.HardwareAddress,
		CardType:        plan.CardType,
	}

	done, err := r.client.DeviceCreate(ctx, deviceCardModel)
	if err != nil || !done {
		resp.Diagnostics.AddError(
			"Error creating Device Card",
			fmt.Sprintf("Error creating Device Card %s: %s", deviceCardModel.DeviceName, err.Error()),
		)
		return
	}

	plan.ID = deviceCardModel.DeviceName
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceCardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceCardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.GetDeviceCard(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Device Card",
			fmt.Sprintf("Could not read Device Card with ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.DeviceName = device.DeviceName
	state.HardwareAddress = device.HardwareAddress
	state.CardType = device.CardType
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceCardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceCardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceCardModel := DeviceCardModel{
		DeviceName:      plan.DeviceName,
		HardwareAddress: plan.HardwareAddress,
		CardType:        plan.CardType,
	}

	_, err := r.client.UpdateDeviceCard(ctx, plan.ID.ValueString(), deviceCardModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Device Card",
			fmt.Sprintf("Could not update Device Card %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	device, err := r.client.GetDeviceCard(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated Device Card",
			fmt.Sprintf("Could not read Device Card %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	state := DeviceCardResourceModel{
		ID:              device.DeviceName,
		DeviceName:      device.DeviceName,
		HardwareAddress: device.HardwareAddress,
		CardType:        device.CardType,
		LastUpdated:     types.StringValue(time.Now().Format(time.RFC850)),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceCardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceCardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDeviceCard(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Device Card",
			fmt.Sprintf("Could not delete Device Card %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *deviceCardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type DeviceCardModel struct {
	DeviceName      types.String `tfsdk:"device_name"`
	HardwareAddress types.String `tfsdk:"hardware_address"`
	CardType        types.String `tfsdk:"card_type"`
}

type DeviceCardResourceModel struct {
	ID              types.String `tfsdk:"id"`
	DeviceName      types.String `tfsdk:"device_name"`
	HardwareAddress types.String `tfsdk:"hardware_address"`
	CardType        types.String `tfsdk:"card_type"`
	LastUpdated     types.String `tfsdk:"last_updated"`
}
