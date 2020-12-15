endpoint "greeter" "http" {
	endpoint = "/"
	method = "POST"
	codec = "json"
}

endpoint "greeter" "grpc" {
	package = "semaphore.greeter"
	service = "Say"
	method = "Hello"
}

flow "greeter" {
	input {
		payload "semaphore.greeter.Request" {}
	}

	resource "user" {
		request "semaphore.greeter.Say" "Hello" {
			name = "{{ input:name }}"
		}
	}

	output {
		payload "semaphore.greeter.Response" {
			msg = "{{ user:msg }}"
			meta = "{{ user:meta }}"
		}
	}
}
