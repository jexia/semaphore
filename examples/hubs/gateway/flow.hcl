endpoint "gateway" "http" {
	endpoint = "/v1/*endpoint"
	method = "GET"
}

proxy "gateway" {
	call "query" {
		request "proto.Users" "ValidateUser" {
		}
	}

	forward "proto.Hub" {
	}
}