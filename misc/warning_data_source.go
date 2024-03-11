package misc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &WarningDataSource{}

func NewWarningDataSource() datasource.DataSource {
	return &WarningDataSource{}
}

// WarningDataSource defines the data source implementation.
type WarningDataSource struct{}

// WarningDataSourceModel describes the data source data model.
type WarningDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Condition types.Bool   `tfsdk:"condition"`
	Summary   types.String `tfsdk:"summary"`
	Details   types.String `tfsdk:"details"`
}

func (d *WarningDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warning"
}

func (d *WarningDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Error data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Error identifier",
				Computed:            true,
			},
			"condition": schema.BoolAttribute{
				MarkdownDescription: "Error condition",
				Required:            true,
			},
			"summary": schema.StringAttribute{
				MarkdownDescription: "Error message summary",
				Required:            true,
			},
			"details": schema.StringAttribute{
				MarkdownDescription: "Error message details",
				Optional:            true,
			},
		},
	}
}

func (d *WarningDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ErrorDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Condition.ValueBool() {
		var details string
		if !data.Details.IsNull() {
			details = data.Details.String()
		}
		resp.Diagnostics.AddWarning(data.Summary.String(), details)
	}
}
