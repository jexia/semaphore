endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
}

service "todo" "http" "json" {
	host = "https://jsonplaceholder.typicode.com"
	schema = "proto.TODO"
}

proxy "todo" {
	call "query" {
		request "todo" "User" {
		}
	}

	forward "todo" "Get" {
	}
}