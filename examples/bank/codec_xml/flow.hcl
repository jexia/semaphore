endpoint "FetchCountry" "http" {
  endpoint = "/"
  method   = "GET"
  codec    = "json"
}

flow "FetchCountry" {
  resource "query" {
    request "com.semaphore.Country" "Get" {}
  }

  output "com.semaphore.Countries" {}
}
