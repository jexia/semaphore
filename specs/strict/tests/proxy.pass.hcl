service "com.maestro" "caller" "http" "json" {
	host = ""
}

proxy "echo" {
	resource "opening" {
		request "caller" "Open" {}
	}

	resource "reference" {
		request "caller" "Open" {
			message = "{{ opening:message }}"
		}
	}

	forward "caller" {}
}