caller "http" {}

service "caller" "http" {
	host = ""
	schema = ""
}

flow "echo" {
    input {
        message = "<string>"
    }

	call "opening" "caller.Open" {
		request {
			message = "{{ input:message }}"
		}
	}
}