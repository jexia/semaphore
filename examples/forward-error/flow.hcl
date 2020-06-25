endpoint "NotFound" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

error "proto.Error" {
	message = "{{ error:message }}"
	status = "{{ error:status }}"
}

flow "NotFound" {
	input "proto.Empty" {}

	resource "query" {
		request "proto.Service" "ThrowError" {
		}
	}

	output "proto.Empty" {}
}
