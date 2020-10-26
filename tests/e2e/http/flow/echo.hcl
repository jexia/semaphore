endpoint "typetest" "http" {
  endpoint = "/json"
  method   = "POST"
  codec    = "json"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  output "semaphore.typetest.Response" {
    echo = "{{ input:data }}"
  }
}
