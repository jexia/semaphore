endpoint "todo" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

proxy "todo" {
	call "query" {
		request "proto.Todo" "User" {
		}
	}

	forward "proto.Forward" {
	}
}