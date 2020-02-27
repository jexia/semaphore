service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
    input {
        ammount = "<int32>"
    }

	call "opening" {
		request "caller" "Open" {
			header {
                Ammount = "{{ input:ammount }}"
            }
		}
	}
}