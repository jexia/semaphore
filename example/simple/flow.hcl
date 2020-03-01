endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
}

service "todo" "http" "json" {
	host = "https://jsonplaceholder.typicode.com"
	schema = "proto.TODO"
}

flow "todo" {
	schema = "proto.TODO.Get"

	call "query" {
		request "todo" "Get" {
		}
	}

	call "user" {
		request "todo" "User" {
		}
	}

	output {
		header {
			username = "{{ user:username }}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}