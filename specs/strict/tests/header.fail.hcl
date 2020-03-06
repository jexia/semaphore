service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
	input "input" {
	}

	call "opening" {
		request "caller" "Open" {
			header {
                Amount = "{{ input:amount }}"
            }
		}
	}
}