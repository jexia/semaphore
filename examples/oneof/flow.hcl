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

  resource "user" {
    request "semaphore.oneof.TestOneOf" "Do" {
      first  = "{{ input:first }}"
      second = "{{ input:second }}"
    }
  }

  output "semaphore.oneof.Response" {
    // msg  = "{{ user:msg }}"  // meta = "{{ user:meta }}"
  }
}
