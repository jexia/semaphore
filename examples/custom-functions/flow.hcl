endpoint "FetchLatestProject" "http" {
	endpoint = "/"
	method	 = "GET"
	codec 	 = "json"
}

error {
	status = 400

	payload "proto.Error" {
		value = "some message"
	}
}

flow "FetchLatestProject" {
	input {
        header = ["Authorization"]

		payload "proto.Query" {}
    }

	error  {
		// TODO: allow reference
		status = 401 // "{{ error:status }}"

		payload "proto.Unauthorized" {
			message = "{{ error:message }}"
		}
	}

	on_error {
		status  = 401
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

	output {
		header {
			Username = "{{ user:username }}"
		}

		payload "proto.Item" {
			id 		  = "{{ query:id }}"
			title 	  = "{{ query:title }}"
			completed = "{{ query:completed }}"
		}
	}
}
