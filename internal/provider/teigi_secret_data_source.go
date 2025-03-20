package landb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &teigiSecretDataSource{}
	_ datasource.DataSourceWithConfigure = &teigiSecretDataSource{}
)

func NewTeigiSecretDataSource() datasource.DataSource {
	return &teigiSecretDataSource{}
}

type teigiSecretDataSource struct {
	client *landb.Client
}

func (d *teigiSecretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teigiSecret"
}

func (d *teigiSecretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostgroup": schema.StringAttribute{
				Required:    true,
				Description: "Hostgroup where the secret is located",
			},
			"key": schema.StringAttribute{
				Optional:    true,
				Description: "Key name which to retrieve",
			},
			"secret": schema.ListAttribute{
				Optional:    true,
				Description: "Secret string retrieved from Teigi",
			},
		},
	}
}

func (d *teigiSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TeigiSecretModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretValue, err := d.client.GetTeigiSecret(state.Hostgroup.ValueString(), state.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Teigi Secret",
			fmt.Sprintf("Could not retrieve secret for hostgroup '%s' and key '%s': %s", state.Hostgroup.ValueString(), state.Key.ValueString(), err.Error()),
		)
		return
	}

	state.Secret = types.StringValue(secretValue)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *teigiSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

type TeigiSecretModel struct {
	Hostgroup types.String `tfsdk:"hostgroup"`
	Key       types.String `tfsdk:"key"`
	Secret    types.String `tfsdk:"secret"`
}
