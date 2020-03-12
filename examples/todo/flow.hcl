endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
}

flow "todo" {
	input "proto.Query" {}

	call "query" {
		request "proto.TODO" "Get" {
		}
	}

	call "user" {
		request "proto.TODO" "User" {
		}
	}

	output "proto.Item" {
		header {
			Username = "{{ user:username }}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
