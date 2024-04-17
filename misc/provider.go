package misc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &kiwiProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &kiwiProvider{}
}

// kiwiProvider is the provider implementation.
type kiwiProvider struct{}

// Metadata returns the provider type name.
func (p *kiwiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "misc"
}

// Schema defines the provider-level schema for configuration data.
func (p *kiwiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform helpers.",
	}
}

// Configure prepares  nothing
func (p *kiwiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Kiwi provider")
}

// DataSources defines the data sources implemented in the provider.
func (p *kiwiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewErrorDataSource,
		NewWarningDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *kiwiProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClaimFromPoolResource,
		NewStatefulListResource,
	}
}
