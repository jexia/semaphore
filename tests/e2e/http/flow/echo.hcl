endpoint "typetest" "http" {
  endpoint = "/json"
  method   = "POST"
  codec    = "json"
}

endpoint "typetest" "http" {
  endpoint = "/xml"
  method   = "POST"
  codec    = "xml"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  output {
    payload "semaphore.typetest.Response" {
      echo = "{{ input:data }}"
    }
  }
}
