endpoint "latest_todo" "graphql" {
	path = "todo.latest"
	name = "LatestTodo"
	base = "query"
}

flow "latest_todo" {
	input {
		payload "com.semaphore.Empty" {}
	}

	resource "query" {
		request "com.semaphore.Todo" "First" {
		}
	}

	output {
		payload "com.semaphore.Item" {
			id = "{{ query:id }}"
			title = "{{ query:title }}"
			completed = "{{ query:completed }}"
		} 
	}
}

endpoint "todo" "graphql" {
	path = "todo.query"
	name = "TodoQuery"
	base = "query"
}

flow "todo" {
	input {
		payload "com.semaphore.Query" {}
	}

	resource "query" {
		request "com.semaphore.Todo" "Get" {
			params {
				id = "{{ input:id }}"
			}
		}
	}

	output {
		payload "com.semaphore.Item" {
			id 		  = "{{ query:id }}"
			userId 	  = "{{ query:userId }}"
			title 	  = "{{ query:title }}"
			completed = "{{ query:completed }}"
		}
	}
}
