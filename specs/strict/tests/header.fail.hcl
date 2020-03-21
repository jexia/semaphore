service "com.maestro" "caller" "http" "json" {
	host = ""
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