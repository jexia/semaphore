endpoint "typetest" "http" {
  endpoint = "/"
  method   = "POST"
  codec    = "json"
}

endpoint "typetest" "grpc" {
  package = "semaphore.typetest"
  service = "Typetest"
  method  = "Run"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  output "semaphore.typetest.Response" {
    oneof = "{{ input:oneof }}"

    object = "{{ input:object }}"

    array = "{{ input:array }}"
  }
}