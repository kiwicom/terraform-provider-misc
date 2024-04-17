package misc

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &statefulList{}
	_ resource.ResourceWithImportState = &statefulList{}
	_ resource.ResourceWithModifyPlan  = &statefulList{}
)

// NewStatefulList is a helper function to simplify the provider implementation.
func NewStatefulListResource() resource.Resource {
	return &statefulList{}
}

// statefulList is the resource implementation.
type statefulList struct{}

// statefulListModel maps the resource schema data.
type statefulListModel struct {
	ID     basetypes.StringValue `tfsdk:"id"`
	Input  basetypes.SetValue    `tfsdk:"input"`
	Output basetypes.SetValue    `tfsdk:"output"`
}

// Metadata returns the data source type name.
func (r *statefulList) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stateful_list"
}

// Schema defines the schema for the data source.
func (r *statefulList) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Stateful list takes items from the input and preserve them in the output. " +
			"The item will always be preserved in the output even if removed from the input. " +
			"Once in, always out!",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Random id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"input": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Set of strings to preserve in the output.",
				Required:    true,
			},
			"output": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Always preserved input. Once in, always out.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *statefulList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statefulListModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(time.Now().Format(time.RFC3339Nano))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *statefulList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *statefulList) update(ctx context.Context, tfplan tfsdk.Plan, tfstate tfsdk.State, diag *diag.Diagnostics) (plan statefulListModel) {
	diags := tfplan.Get(ctx, &plan)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	var state statefulListModel
	tfstate.Get(ctx, &state)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	stateOutput := []string{}
	if !tfstate.Raw.IsNull() {
		diags = state.Output.ElementsAs(ctx, &stateOutput, false)
		diag.Append(diags...)
		if diag.HasError() {
			return
		}
	}

	planInput := []string{}
	diag.Append(plan.Input.ElementsAs(ctx, &planInput, false)...)
	if diag.HasError() {
		return
	}

	for _, i := range planInput {
		if !stringInSlice(i, stateOutput) {
			stateOutput = append(stateOutput, i)
		}
	}

	planOutput, diags := basetypes.NewSetValueFrom(ctx, types.StringType, stateOutput)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	plan.Output = planOutput

	return
}

// ModifyPlan modifies plan in a way to show all the changes before the apply
func (r *statefulList) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// don't modify on delete
	if req.Plan.Raw.IsNull() {
		return
	}

	plan := r.update(ctx, req.Plan, req.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *statefulList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.update(ctx, req.Plan, req.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *statefulList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (r *statefulList) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
