package landb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &egroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &egroupsDataSource{}
)

func NewEGroupDataSource() datasource.DataSource {
	return &egroupsDataSource{}
}

type egroupsDataSource struct {
	client *landb.Client
}

func (d *egroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_egroups"
}

func (d *egroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"query_mails": schema.StringAttribute{
				Optional: true,
			},
			"members": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"mails": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *egroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state egroupsDataSourceModel

	egroups, err := d.client.GetEGroup()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read EGroup",
			err.Error(),
		)
		return
	}

	var memberVals []attr.Value
	for _, m := range egroups.Members {
		memberVals = append(memberVals, types.StringValue(m))
	}

	var mailVals []attr.Value
	for _, m := range egroups.Mails {
		mailVals = append(mailVals, types.StringValue(m))
	}

	for _, egroup := range egroups {
		egroupState := egroupsModel{
			Name:       types.StringValue(egroup.Name),
			QueryMails: types.StringValue(egroup.QueryMails),
			Members:    types.ListValueMust(types.StringType, memberVals),
			Mails:      types.ListValueMust(types.StringType, mailVals),
		}

		state.EGroup = append(state.EGroup, egroupState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *egroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type egroupsDataSourceModel struct {
	EGroup []egroupsModel `tfsdk:"egroups"`
}

type egroupsModel struct {
	Name       types.String `tfsdk:"name"`
	QueryMails types.String `tfsdk:"query_mails"`
	Members    types.List   `tfsdk:"members"`
	Mails      types.List   `tfsdk:"mails"`
}
