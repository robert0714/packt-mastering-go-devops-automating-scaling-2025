terraform {
  required_providers {
    urlshortener = {
      source  = "local/urlshortener" 
      version = ">= 0.0.1" 
    }
  }
}

provider "urlshortener" {
   api_key = "my-secret-api-key"
}

resource "urlshortener" "example" {
  # short_url = "my-url"  # ❌ remove the line - Cannot set value for this attribute as the provider has marked it as read-only.
  long_url  = "https://example.com"
  
  lifecycle {
    ignore_changes = [
      long_url
    ]
  }
}

# 如果需要使用生成的 short_url，可以通過 output 輸出
output "short_url" {
  value = urlshortener.example.short_url
}