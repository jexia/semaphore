caller "http" {}

service "caller" "http" {
	host = ""
	schema = "caller"
	codec = "json"
}

flow "echo" {
    input {
        ammount = "<int32>"
    }

	call "opening" "caller.Open" {
		request {
			header {
                Ammount = "{{ input:ammount }}"
            }
		}
	}
}