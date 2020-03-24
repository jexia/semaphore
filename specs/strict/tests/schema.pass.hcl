service "com.maestro" "caller" "http" "json" {
	host = ""
}

flow "echo" {
  input "input" {
	}

	resource "opening" {
		request "caller" "Open" {
			message = "{{ input:message }}"
		}
	}

	output "output" {
		message = "{{ input:message }}"
	}
}