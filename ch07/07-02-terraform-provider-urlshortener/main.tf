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

# This is the Create and Read test:
resource "urlshortener" "example" {
  long_url = "https://example.com"
}

# This is the Update test:
resource "urlshortener" "example_update" {
  long_url = "https://example2.com"
}