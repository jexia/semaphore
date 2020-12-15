endpoint "FetchLatestProject" "http" {
    endpoint = "/"
    method   = "GET"
    codec    = "json"
}

flow "FetchLatestProject" {
    input {
        header = ["Authorization", "Timestamp"]

        payload "com.semaphore.Query" {}
    }

    resource "query" {
        request "com.semaphore.Todo" "Get" {}
    }

    resource "user" {
      request "com.semaphore.Users" "Get" {}
    }

    output {
      status = 202

      header {
        Username = "{{ user:username }}"
      }

      // TODO: fix panic when message is not defined
      payload "com.semaphore.Item" {
        id        = "{{ query:id }}"
        title     = "{{ query:title }}"
        completed = "{{ query:completed }}"
      }
    }
}
