service "caller" "http" "json" {
	host = ""
}

flow "echo" {
    input "input" {
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