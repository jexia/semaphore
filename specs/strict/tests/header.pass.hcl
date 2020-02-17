caller "http" {}

service "caller" "http" {
	host = ""
	schema = "caller"
	codec = "json"
}

flow "echo" {
    input {
        header {
            Authorization = "<string>"
        }
    }

	call "opening" "caller.Open" {
		request {
			header {
                Authorization = "{{ input.request.header:Authorization }}"
            }
		}
	}
}