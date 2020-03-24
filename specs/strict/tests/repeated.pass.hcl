service "com.maestro" "caller" "http" "json" {
	host = ""
}

flow "echo" {
  	input "input" {
	}

	resource "opening" {
		request "caller" "Open" {
			repeating = "{{ input:repeating }}"
		}
	}
}