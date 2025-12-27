package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider implements framework.Provider
var _ provider.Provider = &urlShortenerProvider{}

// urlShortenerClient is a placeholder for your actual client implementation
type urlShortenerClient struct {
	APIKey     string
	HTTPClient *http.Client
	// simple in-memory store for shortened URLs
	Store     map[string]string
	mu        *sync.Mutex
	storeFile string
}

// Load store from JSON file if configured
func (c *urlShortenerClient) Load() error {
	if c.storeFile == "" {
		return nil
	}
	f, err := os.Open(c.storeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	return decoder.Decode(&c.Store)
}

// Save store to JSON file if configured
func (c *urlShortenerClient) Save() error {
	if c.storeFile == "" {
		return nil
	}
	tmp := c.storeFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c.Store); err != nil {
		f.Close()
		return err
	}
	f.Close()
	return os.Rename(tmp, c.storeFile)
}

// New returns a factory for the provider (used by providerserver).
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &urlShortenerProvider{
			version: version,
		}
	}
}

type urlShortenerProvider struct {
	version string
	apiKey  string
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
			"store_file": schema.StringAttribute{
				Optional:    true,
				Description: "Optional path to a JSON file to persist shortened URLs between runs.",
			},
		},
	}
}

// Resources returns the provider's resources.
func (p *urlShortenerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &URLShortenerResource{} },
	}
}

// DataSources returns the provider's data sources.
func (p *urlShortenerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// provider-level model mapping (used with tfsdk tags)
type urlShortenerProviderModel struct {
	APIKey    types.String `tfsdk:"api_key"`
	StoreFile types.String `tfsdk:"store_file"`
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
		Store:      make(map[string]string),
		mu:         &sync.Mutex{},
	}
	// store file from provider config or env
	if !config.StoreFile.IsNull() {
		client.storeFile = config.StoreFile.ValueString()
	} else if v := os.Getenv("URLSHORTENER_STORE_FILE"); v != "" {
		client.storeFile = v
	}
	// attempt to load existing store
	if err := client.Load(); err != nil {
		resp.Diagnostics.AddError("Failed to load store file", err.Error())
		return
	}
	// Make client available to resources and data sources.
	resp.ResourceData = client
	resp.DataSourceData = client
	p.apiKey = apiKey
}
