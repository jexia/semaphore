service "com.semaphore" "users" {
  transport = "http"
  codec     = "proto"
  host      = "https://users.com"

  options {
    sample = "value"
  }

  method "Add" {
    request  = "proto.Request"
    response = "proto.Response"
  }

  method "Delete" {
    request  = "proto.Request"
    response = "proto.Response"
  }
}

flow "mock" {
  resource "add" {
    request "com.semaphore.users" "Add" {
      repeated "repeated" "input:repeated" {
        key = "input:repeated.value"
      }

      message "nested" {
        key = "value"

        message "nested" {
          key = "value"
        }
      }
    }
  }
}
