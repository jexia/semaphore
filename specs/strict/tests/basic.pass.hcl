service "com.maestro" "caller" "http" "json" {
	host = ""
}

flow "echo" {
	input "input" {
	}

	call "opening" {
		request "caller" "Open" {
			message = "{{ input:message }}"
		}
	}

	output "output" {
		message = "{{ opening:message }}"
	}
}