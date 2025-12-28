terraform {
  required_providers {
    urlshortener = {
      source  = "local/urlshortener" 
      version = "0.1.0" 
    }
  }
}

provider "urlshortener" {
  
}


provider "urlshortener" {
  api_key = "my-secret-api-key"
}
resource "urlshortener_short_url" "example" {
  long_url = "https://example.com"
}