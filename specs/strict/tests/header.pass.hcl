service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
    input {
        header {
            Authorization = "<string>"
        }
    }

	call "opening" {
		request "caller" "Open" {
			header {
                Authorization = "{{ input.header:Authorization }}"
            }
		}
	}
}