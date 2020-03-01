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

	output {
		header {
			example = "{{query:title}}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}