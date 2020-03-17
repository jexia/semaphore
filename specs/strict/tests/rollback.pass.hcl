service "caller" "http" "json" {
	host = ""
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