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
    message = "hello flow"
    // enum = "{{ input:data.enum }}"
    // string = "{{ input:data.string }}"
    // int64 = "{{ input:data.int64 }}"
    dataObject = "{{ input:data.dataObject }}"
  }
}