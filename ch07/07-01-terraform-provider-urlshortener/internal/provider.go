package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider implements framework.Provider
var _ provider.Provider = &urlShortenerProvider{}

// New returns a factory for the provider (used by providerserver).
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &urlShortenerProvider{
			version: version,
		}
	}
}

// urlShortenerClient is a placeholder for your actual client implementation
type urlShortenerClient struct {
	APIKey     string
	HTTPClient *http.Client
}

type urlShortenerProvider struct {
	version string
	apiKey  string
}

// provider-level model mapping (used with tfsdk tags)
type urlShortenerProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

// Metadata returns the provider type name and version.
func (p *urlShortenerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "urlshortener"
	resp.Version = p.version
}

// Schema defines provider-level schema for configuration.
func (p *urlShortenerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Description: "API key for authenticating with the URL Shortener service.",
			},
		},
	}
}

// Resources returns the provider's resources.
func (p *urlShortenerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &urlShortenerResource{} },
	}
}

// DataSources returns the provider's data sources.
func (p *urlShortenerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Configure - build a client and make it available to resources/data-sources
func (p *urlShortenerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Map provider config to the typed model
	var config urlShortenerProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If the value is unknown at configure time, produce a useful attribute
	// diagnostic.
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown API Key",
			"The provider cannot create a URL Shortener client as the `api_key` value is unknown. "+
				"Set it in the provider block or provide the value via the URLSHORTENER_API_KEY environment variable.",
		)
		return
	}
	// Environment fallback: env var wins unless provider config provided one
	apiKey := os.Getenv("URLSHORTENER_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API Key",
			"No API key was provided. Set `api_key` in the provider configuration or set the URLSHORTENER_API_KEY environment variable.",
		)
		return
	}
	// Construct the client
	client := &urlShortenerClient{
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
	// Make client available to resources and data sources.
	resp.ResourceData = client
	resp.DataSourceData = client
	p.apiKey = apiKey
}

func URLResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"short_url": schema.StringAttribute{
				Required:    true,
				Description: "The shortened URL.",
			},
			"long_url": schema.StringAttribute{
				Required:    true,
				Description: "The original long URL.",
			},
		},
	}
}

func CreateShortURL(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data struct {
		ShortURL string `tfsdk:"short_url"`
		LongURL  string `tfsdk:"long_url"`
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	// Call external API to create short URL
	shortURL, err := CreateShortURLAPI(data.LongURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating short URL",
			err.Error(),
		)
		return
	}
	data.ShortURL = shortURL
	resp.State.Set(ctx, &data)
}

// CreateShortURLAPI is a minimal placeholder implementation used by the
// example provider. Replace this with real API calls in a production
// implementation. It returns a deterministic pseudo-short URL for the
// provided long URL.
func CreateShortURLAPI(longURL string) (string, error) {
	if longURL == "" {
		return "", fmt.Errorf("long URL is empty")
	}
	// Use a URL-escaped representation of the long URL as a simple
	// deterministic suffix for the short URL in this example.
	escaped := url.QueryEscape(longURL)
	return "https://short.example/" + escaped, nil
}
