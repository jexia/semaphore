service "com.maestro" "caller" {
	transport = "http"
	codec = "json"
	host = ""

	method "Open" {
		request = "input"
		response = "output"
	}
}

endpoint "echo" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

flow "echo" {
	input "input" {
	}

	resource "opening" {
		request "com.maestro.caller" "Open" {
			message = "{{ input:message }}"
		}
	}

	output "output" {
		message = "{{ opening:message }}"
	}
}