package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the resource implements the resource.Resource interface
var _ resource.Resource = &urlShortenerResource{}

type urlShortenerResource struct {
}

type urlShortenerResourceModel struct {
	ID       types.String `tfsdk:"id"`
	LongURL  types.String `tfsdk:"long_url"`
	ShortURL types.String `tfsdk:"short_url"`
}

func (r *urlShortenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "urlshortener_short_url"
}

func (r *urlShortenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID for the short URL resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"short_url": schema.StringAttribute{
				Computed:    true,
				Description: "The shortened URL.",
			},
			"long_url": schema.StringAttribute{
				Required:    true,
				Description: "The original long URL.",
			},
		},
	}
}

func (r *urlShortenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan urlShortenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// call the shared CreateShortURLAPI helper (example implementation)
	short, err := CreateShortURLAPI(plan.LongURL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}

	plan.ShortURL = types.StringValue(short)
	plan.ID = types.StringValue(short)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *urlShortenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state urlShortenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// No external read for the example provider; assume state is authoritative.
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *urlShortenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan urlShortenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// regenerate short URL on update of long_url in this example
	short, err := CreateShortURLAPI(plan.LongURL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	plan.ShortURL = types.StringValue(short)
	plan.ID = types.StringValue(short)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *urlShortenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For example provider we simply remove the state.
	resp.State.RemoveResource(ctx)
}
