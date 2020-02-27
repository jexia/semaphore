service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
    input {
        message = "<string>"
    }

	call "opening" {
		request "caller" "Open" {
			message = "{{ input:message }}"
		}
	}
}