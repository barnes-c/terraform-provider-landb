package provider

import (
	"context"
	"time"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setDataSourceModel struct {
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

type setDataSource struct {
	client *landb.Client
}

func NewSetDataSource() datasource.DataSource {
	return &setDataSource{}
}

func (d *setDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set"
}

func (d *setDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for retrieving an existing set",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"network_domain": schema.StringAttribute{
				Computed: true,
			},
			"responsible": contactSchemaBlock("Responsible entity for the set"),
			"description": schema.StringAttribute{Computed: true},
			"project_url": schema.StringAttribute{Computed: true},
			"receive_notifications": schema.BoolAttribute{Computed: true},
			"version": schema.Int64Attribute{Computed: true},
			"last_updated": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *setDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*landb.Client); ok {
		d.client = client
	}
}

func (d *setDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data setDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ptr, err := d.client.GetSet(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading set", err.Error())
		return
	}

	data.ID = types.StringValue(ptr.Name)
	data.Type = types.StringValue(ptr.Type)
	data.NetworkDomain = types.StringValue(ptr.NetworkDomain)
	data.Responsible = flattenContactObject(ptr.Responsible)
	data.Description = types.StringValue(ptr.Description)
	data.ProjectURL = types.StringValue(ptr.ProjectURL)
	data.ReceiveNotifications = types.BoolValue(ptr.ReceiveNotifications)
	data.Version = types.Int64Value(int64(ptr.Version))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
