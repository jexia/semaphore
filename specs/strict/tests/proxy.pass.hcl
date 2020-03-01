service "caller" "http" "json" {
	host = ""
	schema = "caller"
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

	forward "caller" "Upload" {}
}