endpoint "greeter" "http" {
	endpoint = "/"
	method = "POST"
	codec = "json"
}

endpoint "greeter" "grpc" {
	package = "maestro.greeter"
	service = "Say"
	method = "Hello"
}

flow "greeter" {
	input "maestro.greeter.Request" {}

	resource "user" {
		request "maestro.greeter.Say" "Hello" {
			name = "{{ input:name }}"
		}
	}

	output "maestro.greeter.Response" {
		msg = "{{ user:msg }}"
		meta = "{{ user:meta }}"
	}
}
