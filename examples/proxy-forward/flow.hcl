endpoint "todo" "http" "json" {
	endpoint = "/"
	method = "GET"
}

proxy "todo" {
	call "query" {
		request "proto.Todo" "User" {
		}
	}

	forward "proto.Forward" {
	}
}