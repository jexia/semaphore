endpoint "FetchLatestProject" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

flow "FetchLatestProject" {
	input "proto.Query" {}

	resource "query" {
		request "proto.Service" "GetTodo" {
		}
	}

	resource "user" {
		request "proto.Service" "GetUser" {
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
