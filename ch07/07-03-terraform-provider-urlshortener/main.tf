terraform {
  required_providers {
    urlshortener = {
      source  = "local/urlshortener" 
      version = "0.0.1"
    }
  }
}

provider "urlshortener" {
   api_key = "my-secret-api-key"
}

resource "url_shortener" "example" {
  short_url = "my-url"
  long_url  = "https://example.com"
  lifecycle {
    prevent_destroy = true
  }
}