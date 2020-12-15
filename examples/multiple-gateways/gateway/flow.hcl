endpoint "gateway" "http" {
	endpoint = "/v1/*endpoint"
	method 	 = "GET"
}

proxy "gateway" {
	resource "query" {
		request "com.semaphore.Users" "ValidateUser" {
		}
	}

	forward "com.semaphore.Hub" {
		header {
			User = "{{ query:name }}"
		}
	}
}