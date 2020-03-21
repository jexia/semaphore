endpoint "latest_todo" "graphql" {
}

flow "latest_todo" {
	input "proto.Empty" {}

	call "query" {
		request "proto.TODO" "First" {
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}

endpoint "todo" "graphql" {
}

flow "todo" {
	input "proto.Query" {}

	call "query" {
		request "proto.TODO" "Get" {
			id = "{{ input:id }}"
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
