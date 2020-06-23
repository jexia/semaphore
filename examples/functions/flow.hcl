endpoint "FetchLatestProject" "http" {
	endpoint = "/"
	method = "GET"
	codec = "json"
}

error "placeholder.Error" {
	value = "some message"
	status = 400
}

flow "FetchLatestProject" {
	input "placeholder.Query" {
        header = ["Authorization"]
    }

	error "placeholder.Unauthorized" {
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
		request "placeholder.Service" "GetTodo" {
		}
	}

	resource "user" {
		request "placeholder.Service" "GetUser" {
		}
	}

	output "placeholder.Item" {
		header {
			Username = "{{ user:username }}"
		}

		id = "{{ query:id }}"
		title = "{{ query:title }}"
		completed = "{{ query:completed }}"
	}
}
