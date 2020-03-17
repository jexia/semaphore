service "caller" "http" "json" {
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

	call "opening" {
		request "caller" "Open" {
			message = "{{ input:message }}"
		}

		rollback "caller" "Open" {
			message = "{{ opening:message }}"
		}
	}

	output "output" {
		message = "{{ opening:message }}"
	}
}