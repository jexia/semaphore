endpoint "oneof" "http" {
  endpoint = "/"
  method   = "POST"
  codec    = "json"
}

endpoint "oneof" "grpc" {
  package = "semaphore.oneof"
  service = "TestOneOf"
  method  = "Do"
}

flow "oneof" {
  input "semaphore.oneof.Request" {}

  output "semaphore.oneof.Response" {
    // response = "{{ sprintf('%s %json %s %json', input:first, input:second, input:third, input:fourth) }}"
  }
}
