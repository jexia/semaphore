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
    // oneof = "{{ input:oneof }}"

    // oneof = {
    //   empty = "{{ input:oneof.empty }}"

    //   single = "{{ input:oneof.single }}"

    //   plural = "{{ input:oneof.plural }}"
    // }

    object = "{{ input:object }}"

    array = "{{ input:array }}"
  }
}