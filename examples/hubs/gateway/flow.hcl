endpoint "gateway" "http" {
	endpoint = "/v1/*endpoint"
	method = "GET"
}

proxy "gateway" {
	resource "query" {
		request "proto.Users" "ValidateUser" {
		}
	}

	forward "proto.Hub" {
		header {
			User = "{{ query:name }}"
		}
	}
}