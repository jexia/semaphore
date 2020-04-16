service "com.maestro" "caller" {
	transport = "http"
	codec = "json"
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
		request "com.maestro.caller" "Open" {
			header {
                Authorization = "{{ input.header:Authorization }}"
            }
		}
	}
}