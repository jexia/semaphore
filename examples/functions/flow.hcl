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

	error "placeholder.Error" {
		header {
			X-Request = "{{ error.params:id }}"
		}

		value = "{{ error:message }}"
		status = "{{ error:status }}"
	}

	on_error {
		schema = "placeholder.Error"
		status = 500
		message = "on error message"

		params {
			id = "abc"
		}
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

		error "placeholder.Error" {
			header {
				X-Request = "{{ error.params:id }}"
			}

			value = "{{ error:message }}"
			status = "{{ error:status }}"
		}

		on_error {
			schema = "placeholder.Error"
			status = 500
			message = "on  blablaerror message"

			params {
				id = "abc"
			}
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
