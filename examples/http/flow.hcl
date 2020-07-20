endpoint "FetchLatestProject" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

flow "FetchLatestProject" {
	input "com.semaphore.Query" {}

	resource "query" {
		request "com.semaphore.Service" "GetTodo" {
		}
	}

	resource "user" {
		request "com.semaphore.Service" "GetUser" {
		}
	}

	output "com.semaphore.Item" {
		header {
			Username = "{{ user:username }}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
