// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package provider

import (
	"context"
	"os"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &landbProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &landbProvider{
			version: version,
		}
	}
}

type LandbModel struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Audience     types.String `tfsdk:"audience"`
}

type landbProvider struct {
	version string
}

func (p *landbProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "landb"
	resp.Version = p.version
}

func (p *landbProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional: true,
			},
			"audience": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *landbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring LanDB client")

	var config LandbModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("LANDB_ENDPOINT")
	client_id := os.Getenv("LANDB_SSO_CLIENT_ID")
	client_secret := os.Getenv("LANDB_SSO_CLIENT_SECRET")
	audience := os.Getenv("LANDB_SSO_AUDIENCE")

	ctx = tflog.SetField(ctx, "endpoint", endpoint)
	ctx = tflog.SetField(ctx, "client_id", client_id)
	ctx = tflog.SetField(ctx, "client_secret", client_secret)
	ctx = tflog.SetField(ctx, "audience", audience)

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Invalid LanDB API endpoint",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_ENDPOINT environment variable.",
		)
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Invalid LanDB LDAP client_id",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API client_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_SSO_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Invalid LanDB API client_secret",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API client_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_SSO_CLIENT_SECRET environment variable.",
		)
	}

	if config.Audience.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("audience"),
			"Unknown  audience",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API audience. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_SSO_AUDIENCE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ClientID.IsNull() {
		client_id = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	if !config.Audience.IsNull() {
		audience = config.Audience.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing LanDB API endpoint",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the endpoint. "+
				"Set the endpoint value in the configuration or use the LANDB_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing CERN SSO client ID",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the client_id. "+
				"Set the client_id value in the configuration or use the LANDB_SSO_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing CERN SSO client secret",
			"The provider cannot fetch a authentication token from the CERN SSO application as there is a missing or empty value for the client_secret. "+
				"Set the client_secret value in the configuration or use the LANDB_SSO_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if audience == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("audience"),
			"Missing CERN SSO audience",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the LanDB API audience. "+
				"Set the audience value in the configuration or use the LANDB_SSO_AUDIENCE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "endpoint", endpoint)
	ctx = tflog.SetField(ctx, "client_id", client_id)
	ctx = tflog.SetField(ctx, "client_secret", client_secret)
	ctx = tflog.SetField(ctx, "audience", audience)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "password")

	tflog.Debug(ctx, "Creating LanDB client")

	client, err := landb.NewClient(endpoint, client_id, client_secret, audience)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create LanDB API Client",
			"An unexpected error occurred when creating the LanDB API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"LanDB Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured LanDB client", map[string]any{"success": true})
}

func (p *landbProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewSetAttachmentResource,
		NewSetResource,
	}
}

func (p *landbProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDeviceDataSource,
		NewSetDataSource,
	}
}
