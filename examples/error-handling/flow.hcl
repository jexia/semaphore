endpoint "GlobalHandleError" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

error "proto.Error" {
	message = "{{ error:message }}"
	status = "{{ error:status }}"
}

flow "GlobalHandleError" {
	input "proto.Empty" {}

	on_error {
		status = 501
		message = "Service unavailable"
	}

	resource "query" {
		request "proto.Service" "FetchUser" {
		}

		expect_status = 404

		error "proto.Error" {
			message = "{{ query.error:name }}"
			status = "{{ error:status }}"
		}

		on_error {
			schema = "proto.User"
			message = "Unauthorized"

			params {
				reason = "{{ error:message }}"
			}
		}
	}

	output "proto.Empty" {}
}

endpoint "FlowHandleError" "http" {
	endpoint = "/flow"
	method = "GET"
	codec = "json"
}

flow "FlowHandleError" {
	input "proto.Empty" {}

	error "proto.Error" {
		message = "{{ error:message }}"
		status = "{{ error:status }}"
	}

	on_error {
		status = 401
		message = "flow error message"
	}

	resource "query" {
		request "proto.Service" "ThrowError" {
		}
	}

	output "proto.Empty" {}
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

		error "proto.Error" {
			message = "{{ error:message }}"
			status = "{{ error:status }}"
		}

		on_error {
			status = 401
			message = "node error message"
		}
	}

	output "proto.Empty" {}
}