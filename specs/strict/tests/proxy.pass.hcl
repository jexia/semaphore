service "com.maestro" "caller" "http" "json" {
	host = ""
}

proxy "echo" {
	call "opening" {
		request "caller" "Open" {}
	}

	call "reference" {
		request "caller" "Open" {
			message = "{{ opening:message }}"
		}
	}

	forward "caller" {}
}