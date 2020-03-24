endpoint "latest_todo" "graphql" {
	path = "latest.todo"
	name = "LatestTodo"
	base = "query"
}

flow "latest_todo" {
	input "proto.Empty" {}

	resource "query" {
		request "proto.Todo" "First" {
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}

endpoint "todo" "graphql" {
	path = "todo"
	name = "Todo"
	base = "query"
}

flow "todo" {
	input "proto.Query" {}

	resource "query" {
		request "proto.Todo" "Get" {
			id = "{{ input:id }}"
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
