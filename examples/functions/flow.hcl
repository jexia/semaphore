endpoint "FetchLatestProject" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

error "proto.Error" {
	value = "some message"
	status = 400
}

flow "FetchLatestProject" {
	input "proto.Query" {
        header = ["Authorization"]
    }

	error "proto.Unauthorized" {
		message = "{{ error:message }}"
		status = "{{ error:status }}"
	}

	on_error {
		status = 401
		message = "on error message"
	}

    before {
        resources {
            token = "{{ jwt(input.header:Authorization) }}"
        }
    }

	resource "query" {
		request "proto.Service" "GetTodo" {
		}
	}

	resource "user" {
		request "proto.Service" "GetUser" {
		}
	}

	output "proto.Item" {
		header {
			Username = "{{ user:username }}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
