endpoint "greeter" "http" {
	endpoint = "/"
	method = "POST"
	codec = "json"
}

endpoint "greeter" "graphql" {
	path = "greeter"
}

flow "greeter" {
	input "go.micro.srv.greeter.Request" {}

	call "user" {
		request "go.micro.srv.greeter.Say" "Hello" {
			name = "{{ input:name }}"
		}
	}

	output "go.micro.srv.greeter.Response" {
		msg = "{{ user:msg }}"
	}
}
