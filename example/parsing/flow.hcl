endpoint "logger" "http" "json" {
	endpoint = "/"
	method = "POST"
}

service "logger" "http" "json" {
	host = "https://jsonplaceholder.typicode.com"
	schema = "proto.Logger"
}

flow "logger" {
	schema = "proto.Logger.Append"

	call "logging" {
		request "logger" "Append" {
			message = "{{ input:message }}"
		}
	}

	output {
        id = "{{ logging:id }}"
    }
}