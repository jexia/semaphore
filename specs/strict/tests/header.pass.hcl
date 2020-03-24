service "com.maestro" "caller" "http" "json" {
	host = ""
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