service "com.maestro" "caller" {
	transport = "http"
	codec = "json"
	host = ""
}

flow "echo" {
	input "com.input" {
	}

	resource "opening" {
		request "caller" "Open" {
			message = "{{ input:message }}"
		}
	}
}