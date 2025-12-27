package provider

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type URLShortenerResource struct{}

type URLShortenerResourceModel struct {
	ID       types.String `tfsdk:"id"`
	LongURL  types.String `tfsdk:"long_url"`
	ShortURL types.String `tfsdk:"short_url"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (r *URLShortenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "urlshortener"
}

func (r *URLShortenerResource) Schema(
	ctx context.Context, request resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the short URL.",
			},
			"long_url": resourceschema.StringAttribute{
				Required:    true,
				Description: "The original long URL.",
			},
			"short_url": resourceschema.StringAttribute{
				Computed:    true,
				Description: "The generated short URL.",
			},
		},
	}
}

func (r *URLShortenerResource) Create(
	ctx context.Context, request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var data URLShortenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	// Generate a short URL and ID
	data.ShortURL = types.StringValue(generateShortURL())
	data.ID = types.StringValue(strconv.Itoa(rand.Intn(1000000)))
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func generateShortURL() string {
	return "short.ly/" + strconv.Itoa(rand.Intn(10000))
}

func (r *URLShortenerResource) Read(
	ctx context.Context, req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read the current state
	var state URLShortenerResourceModel

	// Get the current state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *URLShortenerResource) Update(
	ctx context.Context, req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve the current state and plan
	var state URLShortenerResourceModel
	err := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(err...)
	if resp.Diagnostics.HasError() {
		return
	}
	var plan URLShortenerResourceModel
	err = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(err...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the long URL changed, we'll update the short URL
	if plan.LongURL.ValueString() != state.LongURL.ValueString() {
		state.LongURL = plan.LongURL
		state.ShortURL = types.StringValue(generateShortURL())
	}
	// Save the updated state
	err = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(err...)

}

func (r *URLShortenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Remove the resource from state
	resp.State.RemoveResource(ctx)
}
