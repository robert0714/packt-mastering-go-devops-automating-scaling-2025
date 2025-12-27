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

type URLShortenerResource struct {
	client *urlShortenerClient
}

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
	// Ensure provider client available
	if r.client == nil {
		response.Diagnostics.AddError("Provider Client Missing", "Provider client is not configured")
		return
	}

	// Generate a short URL and ID and persist in provider store
	data.ShortURL = types.StringValue(generateShortURL())
	id := strconv.Itoa(rand.Intn(1000000))
	data.ID = types.StringValue(id)

	r.client.mu.Lock()
	r.client.Store[data.ShortURL.ValueString()] = data.LongURL.ValueString()
	r.client.mu.Unlock()

	// persist
	if err := r.client.Save(); err != nil {
		response.Diagnostics.AddError("Failed to save store", err.Error())
		return
	}

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
	// For import or fresh read, fetch long URL from provider store if available
	if r.client != nil {
		r.client.mu.Lock()
		long, ok := r.client.Store[state.ShortURL.ValueString()]
		r.client.mu.Unlock()
		if ok {
			state.LongURL = types.StringValue(long)
		}
	}

	// Save the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *URLShortenerResource) Update(
	ctx context.Context, req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var state URLShortenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var plan URLShortenerResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.LongURL.ValueString() != state.LongURL.ValueString() {
		state.LongURL = plan.LongURL
		// Keep same short URL, just update mapping
		if r.client != nil {
			r.client.mu.Lock()
			r.client.Store[state.ShortURL.ValueString()] = state.LongURL.ValueString()
			r.client.mu.Unlock()
			if err := r.client.Save(); err != nil {
				resp.Diagnostics.AddError("Failed to save store", err.Error())
				return
			}
		}
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *URLShortenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Remove the resource from provider store and state
	var state URLShortenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.client != nil {
		r.client.mu.Lock()
		delete(r.client.Store, state.ShortURL.ValueString())
		r.client.mu.Unlock()
		if err := r.client.Save(); err != nil {
			resp.Diagnostics.AddError("Failed to save store", err.Error())
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

func (r *URLShortenerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*urlShortenerClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Provider data is not a *urlShortenerClient")
		return
	}
	r.client = client
}

// The provider uses the terraform-plugin-framework. Keep resource implementation
// consistent with the framework API. Any SDK v2 helper/schema-based code was
// removed because the rest of the provider is written against the framework.
