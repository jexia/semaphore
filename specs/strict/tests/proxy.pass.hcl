service "caller" "http" {
	host = ""
	schema = "caller"
	codec = "json"
}

proxy "echo" {
	call "opening" "caller.Open" {
		request {}
	}

	call "reference" "caller.Open" {
		request {
			message = "{{ opening:message }}"
		}
	}

	forward "caller.Upload" {}
}