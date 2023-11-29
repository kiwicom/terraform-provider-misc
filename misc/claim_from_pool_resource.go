package kiwi

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
	_ resource.Resource                = &claimFromPool{}
	_ resource.ResourceWithConfigure   = &claimFromPool{}
	_ resource.ResourceWithImportState = &claimFromPool{}
	_ resource.ResourceWithModifyPlan  = &claimFromPool{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewClaimFromPoolResource() resource.Resource {
	return &claimFromPool{}
}

// claimFromPool is the resource implementation.
type claimFromPool struct{}

// claimgFromPoolModel maps the resource schema data.
type claimgFromPoolModel struct {
	ID       types.String       `tfsdk:"id"`
	Pool     []string           `tfsdk:"pool"`
	Claimers []string           `tfsdk:"claimers"`
	Output   basetypes.MapValue `tfsdk:"output"`
}

// Metadata returns the data source type name.
func (r *claimFromPool) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_claim_from_pool"
}

// Schema defines the schema for the data source.
func (r *claimFromPool) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages pool claimers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Random id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Set of items in the pool claimers will claim. Duplicates are removed.",
				Required:    true,
			},
			"claimers": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "List of claimers. Duplicate are removed.",
				Required:    true,
			},
			"output": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "Map of claimed items from the pool (claimer => pool item)",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func (r *claimFromPool) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan claimgFromPoolModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(plan.Claimers) > len(plan.Pool) {
		resp.Diagnostics.AddError("Number of claimers shouldn't be higher than number of items in the pool",
			"")
		return
	}
}

// Configure adds the provider configured client to the data source.
func (r *claimFromPool) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
}

// Create creates the resource and sets the initial Terraform state.
func (r *claimFromPool) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan claimgFromPoolModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(time.Now().Format(time.RFC3339Nano))

	m := map[string]string{}

	for i, c := range plan.Claimers {
		m[c] = plan.Pool[i]
	}

	mv, diags := basetypes.NewMapValueFrom(ctx, types.StringType, m)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Output = mv

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *claimFromPool) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *claimFromPool) update(ctx context.Context, tfplan tfsdk.Plan, tfstate tfsdk.State, diag *diag.Diagnostics) (plan claimgFromPoolModel) {
	diags := tfplan.Get(ctx, &plan)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	var state claimgFromPoolModel
	tfstate.Get(ctx, &state)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	stateOutput := map[string]string{}
	if !tfstate.Raw.IsNull() {
		diags = state.Output.ElementsAs(ctx, &stateOutput, false)
		diag.Append(diags...)
		if diag.HasError() {
			return
		}
	}

	freePool := make([]string, len(plan.Pool))
	copy(freePool, plan.Pool)
	notYetClaimers := make([]string, len(plan.Claimers))
	copy(notYetClaimers, plan.Claimers)

	for c, p := range stateOutput {
		if !stringInSlice(c, plan.Claimers) || !stringInSlice(p, plan.Pool) {
			delete(stateOutput, c)
		}
	}

	for c, p := range stateOutput {
		notYetClaimers = deleteFromSlice(notYetClaimers, c)
		freePool = deleteFromSlice(freePool, p)
	}

	for k, c := range notYetClaimers {
		stateOutput[c] = freePool[k]
	}

	mv, diags := basetypes.NewMapValueFrom(ctx, types.StringType, stateOutput)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	plan.Output = mv

	return
}

// ModifyPlan modifies plan in a way to show all the changes before the apply
func (r *claimFromPool) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

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
func (r *claimFromPool) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
func (r *claimFromPool) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (r *claimFromPool) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
