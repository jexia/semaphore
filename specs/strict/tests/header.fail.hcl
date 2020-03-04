service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
    input {
        amount = "<int32>"
    }

	call "opening" {
		request "caller" "Open" {
			header {
                Amount = "{{ input:amount }}"
            }
		}
	}
}