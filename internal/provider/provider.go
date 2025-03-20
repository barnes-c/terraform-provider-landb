package landb

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &landb{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &landb{
			version: version,
		}
	}
}

type landb struct {
	version string
}

func (p *landb) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "LanDB"
	resp.Version = p.version
}

func (p *landb) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ldap_server": schema.StringAttribute{
				Optional: true,
			},
			"landb_endpoint": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}
func (p *landb) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring LanDB client")

	var config LandbModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ldap_server := os.Getenv("LANDB_LDAP_SERVER")
	landb_endpoint := os.Getenv("LANDB_ENDPOINT")
	username := os.Getenv("LANDB_USERNAME")
	password := os.Getenv("LANDB_USERNAME")

	ctx = tflog.SetField(ctx, "ldap_server", ldap_server)
	ctx = tflog.SetField(ctx, "landb_endpoint", landb_endpoint)
	ctx = tflog.SetField(ctx, "username", username)
	ctx = tflog.SetField(ctx, "password", password)

	if config.LdapServer.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ldap_server"),
			"Unknown LanDB LDAP server",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API ldap_server. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LDAP_SERVER environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("landb_endpoint"),
			"Unknown LanDB API Endpoint",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_USERNAME environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown LanDB API Username",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown LanDB API Password",
			"The provider cannot create the LanDB API client as there is an unknown configuration value for the LanDB API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LANDB_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.LdapServer.IsNull() {
		ldap_server = config.LdapServer.ValueString()
	}

	if !config.LandbEndpoint.IsNull() {
		landb_endpoint = config.LandbEndpoint.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if ldap_server == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ldap_server"),
			"Missing LanDB LDAP server",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the LanDB LDAP server. "+
				"Set the host value in the configuration or use the LANDB_LDAP_SERVER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if landb_endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("landb_endpoint"),
			"Missing LanDB API Endpoint",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the LanDB API endpoint. "+
				"Set the landb_endpoint value in the configuration or use the LANDB_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing LanDB API Username",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the LanDB API username. "+
				"Set the username value in the configuration or use the LANDB_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing LanDB API Password",
			"The provider cannot create the LanDB API client as there is a missing or empty value for the LanDB API password. "+
				"Set the password value in the configuration or use the LANDB_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ldap_server", ldap_server)
	ctx = tflog.SetField(ctx, "landb_endpoint", landb_endpoint)
	ctx = tflog.SetField(ctx, "username", username)
	ctx = tflog.SetField(ctx, "password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "password")

	tflog.Debug(ctx, "Creating LanDB client")

	client, err := landb.NewClient(&ldap_server, &landb_endpoint, &username, &password)
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

func (p *landb) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEGroupDataSource,
		NewTeigiSecretDataSource,
	}
}

func (p *landb) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewDeviceCardResource,
		NewDeviceInterfaceResource,
	}
}

type LandbModel struct {
	LdapServer    types.String `tfsdk:"ldap_server"`
	LandbEndpoint types.String `tfsdk:"landb_endpoint"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
}
