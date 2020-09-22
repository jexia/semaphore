endpoint "FetchCountry" "http" {
  endpoint = "/"
  method   = "GET"
  codec    = "json"
}

flow "FetchCountry" {
  resource "query" {
    request "com.semaphore.WorldBank" "GetCountries" {}
  }

  output "com.semaphore.CountriesResponse" {
    countries = "{{ query:country }}"
  }
}
