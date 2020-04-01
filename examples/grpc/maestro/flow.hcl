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

	resources {
		value = "{{ input:name }}"
	}

	resources {
		else = "{{ add(value:.) }}"
	}

	resource "user" {
		request "maestro.greeter.Say" "Hello" {
			name = "{{ else: }}"
		}
	}

	output "maestro.greeter.Response" {
		msg = "{{ user:msg }}"
	}
}
