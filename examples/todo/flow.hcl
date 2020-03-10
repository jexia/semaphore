endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
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

endpoint "second" "graphql" "json" {
}

flow "second" {
	input "proto.Query" {}

	call "query" {
		request "proto.TODO" "Second" {
		}
	}

	output "proto.Item" {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}