endpoint "FetchLatestProject" "http" {
  endpoint = "/"
  method   = "GET"
  codec    = "json"
}

flow "FetchLatestProject" {
  input "com.semaphore.Query" {
    header = ["Authorization", "Timestamp"]
  }

  resource "query" {
    request "com.semaphore.Todo" "Get" {}
  }

  resource "user" {
    request "com.semaphore.Users" "Get" {}
  }

  output "com.semaphore.Item" {
    header {
      Username = "{{ user:username }}"
    }

    id        = "{{ query:id }}"
    title     = "{{ query:title }}"
    completed = "{{ query:completed }}"
  }
}
