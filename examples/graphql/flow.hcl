endpoint "latest" "graphql" "json" {
}

flow "latest" {
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

endpoint "todo" "graphql" "json" {
}

flow "todo" {
	input "proto.Query" {}

	call "query" {
		request "proto.TODO" "Get" {
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
