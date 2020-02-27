service "caller" "http" "json" {
	host = ""
	schema = "caller"
}

flow "echo" {
    schema = "caller.Open"

	call "opening" {
		request "caller" "Open" {
			repeating = "{{ input:repeating }}"
		}
	}
}