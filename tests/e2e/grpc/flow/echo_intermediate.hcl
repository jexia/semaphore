endpoint "typetest" "grpc" {
	package = "semaphore.typetest"
	service = "Typetest"
	method  = "Run"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  resource "echo" {
    request "semaphore.typetest.External" "Post" {
      enum    = "{{ input:data.enum }}"
      string  = "{{ input:data.string }}"
      integer = "{{ input:data.integer }}"
      double  = "{{ input:data.double }}"
      numbers = "{{ input:data.numbers }}"
    }
  }

  output {
    payload "semaphore.typetest.Response" {
      echo = "{{ echo:. }}"
    }
  }
}

