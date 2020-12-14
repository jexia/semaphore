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
	input "semaphore.greeter.Request" {}

	resource "user" {
		request "semaphore.greeter.Say" "Hello" {
			name = "{{ input:name }}"
		}
	}

	on_error {
		schema = "semaphore.greeter.Error"
		status = 401
		cause  = "flow error message"

		// params {
		// 	prop = ""
		// }
	}

	output {
		status = 202

		payload "semaphore.greeter.Response" {
			msg  = "{{ user:msg }}"
			meta = "{{ user:meta }}"
		}
	}
}
