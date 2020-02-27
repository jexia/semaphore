endpoint "logger" "http" "json" {
	endpoint = "/"
	method = "POST"
}

service "logger" "http" "json" {
	host = "https://my-json-server.typicode.com"
	schema = "proto.Logger"
	method = "GET"
}

flow "logger" {
	schema = "proto.Logger.Call"

	call "logging" {
		request "logger" "Append" {
			message = "{{ input:message }}"
		}
	}

	output {
        posts = "{{ logging:posts }}"
    }
}