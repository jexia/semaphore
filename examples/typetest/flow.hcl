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

    // Is this correct? (since only one of the fields can be set) Or we should allow 
    // references to entire 'oneof' only?
    oneof "oneof" {
      empty = "{{ input:oneof.empty }}"

      single = "{{ input:oneof.single }}"

      plural = "{{ input:oneof.plural }}"
    }

    array = "{{ input:array }}"
  }
}