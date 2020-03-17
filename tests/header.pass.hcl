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
        header = ["Authorization"]
    }

	call "opening" {
		request "caller" "Open" {
			header {
                Authorization = "{{ input.header:Authorization }}"
            }
		}
	}
}