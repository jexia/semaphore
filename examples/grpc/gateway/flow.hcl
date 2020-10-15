endpoint "greeter" "http" {
  endpoint = "/"
  method   = "POST"
  codec    = "json"
}

endpoint "greeter" "grpc" {
  package = "semaphore.greeter"
  service = "Say"
  method  = "Hello"
}

flow "greeter" {
  input "semaphore.greeter.Request" {}

  resource "user" {
    request "semaphore.greeter.Say" "Hello" {
      name = "{{ input:name }}"
    }
  }

  output "semaphore.greeter.Response" {
    // TODO: here is a trouble!!!!!!!  // get outputs from the flow (not from the specs)  // msg = "{{ user:msg }}"  // meta = "{{ user:meta }}"
  }
}
