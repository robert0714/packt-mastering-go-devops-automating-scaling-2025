terraform {
  required_providers {
    urlshortener = {
      source  = "local/urlshortener" 
      version = ">= 0.0.1" # Accept version 0.0.1 or higher
    }
  }
}

provider "urlshortener" {
  api_key = "my-secret-api-key"
}
resource "urlshortener_short_url" "example" {
  long_url = "https://example.com"
}
