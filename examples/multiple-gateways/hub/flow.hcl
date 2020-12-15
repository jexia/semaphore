endpoint "user" "http" {
	endpoint = "/v1/user"
	method 	 = "GET"
	codec 	 = "json"
}

flow "user" {
	resource "query" {
		request "com.semaphore.Users" "GetUser" {
		}
	}

	output {
		payload "com.semaphore.User" {
			id = "{{ query:id }}"
			name = "{{ query:name }}"
			username = "{{ query:username }}"
		}
    }
}