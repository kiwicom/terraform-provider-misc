package misc

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ErrorDataSource{}

func NewErrorDataSource() datasource.DataSource {
	return &ErrorDataSource{}
}

// ErrorDataSource defines the data source implementation.
type ErrorDataSource struct{}

// ErrorDataSourceModel describes the data source data model.
type ErrorDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Condition types.Bool   `tfsdk:"condition"`
	Summary   types.String `tfsdk:"summary"`
	Details   types.String `tfsdk:"details"`
}

func (d *ErrorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_error"
}

func (d *ErrorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *ErrorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ErrorDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Condition.ValueBool() {
		var details string
		if !data.Details.IsNull() {
			details = data.Details.String()
		}
		resp.Diagnostics.AddError(data.Summary.String(), details)
	}

	types.StringValue(time.Now().Format(time.RFC3339Nano))
}
