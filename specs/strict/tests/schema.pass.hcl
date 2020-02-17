caller "http" {}

service "caller" "http" {
	host = ""
	schema = "caller"
	codec = "json"
}

flow "echo" {
    schema = "caller.Open"

	call "opening" "caller.Open" {
		request {
			message = "{{ input:message }}"
		}
	}

	output {
		message = "{{ input:message }}"
	}
}