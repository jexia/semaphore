caller "http" {
	base = "/v1"
}

service "logger" "http" {
	host = "logger.local"
	schema = "proto.Logger"
    codec = "proto"
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