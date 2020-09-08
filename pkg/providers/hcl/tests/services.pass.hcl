service "com.semaphore" "auth" {
  transport = "http"
  codec     = "json"
  host      = "https://auth.com"
}

service "com.semaphore" "auth" {}

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
