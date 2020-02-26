endpoint "logger" "http" {
	codec = "json"
}

service "logger" "http" {
	host = "logger.local"
	schema = "proto.Logger"
    codec = "proto"

	options {
		endpoint = "/"
		method = "GET"
	}
}

flow "logger" {
	schema = "proto.Logger.Append"

	call "logging" "logger.Append" {
		request {
			message = "{{ input:message }}"
		}
	}

	output {
        id = "{{ logging:id }}"
    }
}