endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
}

service "todo" "http" "json" {
	host = "https://jsonplaceholder.typicode.com"
	schema = "proto.TODO"
}

flow "todo" {
	call "query" {
		request "todo" "Get" {
		}
	}

	output {
		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}