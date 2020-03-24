service "com.maestro" "caller" "http" "json" {
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

	resource "opening" {
		request "caller" "Open" {
			header {
                Authorization = "{{ input.header:Authorization }}"
            }
		}
	}
}