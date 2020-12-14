endpoint "GlobalHandleError" "graphql" {
	path = "greeter.global"
	name = "GlobalError"
	base = "query"
}

endpoint "GlobalHandleError" "grpc" {
	package = "semaphore.greeter"
	service = "Error"
	method = "Global"
}

endpoint "GlobalHandleError" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

error  {
	status = 400 // "{{ error:status }}"
	
	payload "proto.Error" {
		message = "{{ error:message }}"
	}
}

flow "GlobalHandleError" {
	input "proto.Empty" {}

	resource "query" {
		request "proto.Service" "ThrowError" {
		}

		on_error {
			status = 401
			message = "global error message"
		}
	}

	output {
		// TODO: omit empty payload
		payload "proto.Empty" {}
	}
}

endpoint "FlowHandleError" "http" {
	endpoint = "/flow"
	method = "GET"
	codec = "json"
}

flow "FlowHandleError" {
	input "proto.Empty" {}

	error {
		status = 401 // "{{ error:status }}"

		payload "proto.Error" {
			message = "{{ error:message }}"
		}
	}

	on_error {
		status = 402
		message = "flow error message"
	}

	resource "query" {
		request "proto.Service" "ThrowError" {
		}
	}

	output {
		// TODO: omit empty payload
		payload "proto.Empty" {}
	}
}

endpoint "NodeHandleError" "http" {
	endpoint = "/flow/node"
	method = "GET"
	codec = "json"
}

flow "NodeHandleError" {
	input "proto.Empty" {}

	resource "query" {
		request "proto.Service" "ThrowError" {
		}

		error {
			status = 400

			payload "proto.Error" {
				message = "{{ error:message }}"
			}
		}

		on_error {
			status = 403
			message = "node error message"
		}
	}

	output {
		// TODO: omit empty payload
		payload "proto.Empty" {}
	}
}