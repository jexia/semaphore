endpoint "typetest" "grpc" {
	package = "semaphore.typetest"
	service = "Typetest"
	method = "Run"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  output "semaphore.typetest.Response" {
    echo = "{{ input:data }}"
  }
}
