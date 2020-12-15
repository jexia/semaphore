protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
  endpoint = "/user/:id"
  method = "post"
}

flow "CreateUser" {
  input "com.semaphore.User" {}

  output {
    payload "com.semaphore.User" {}
  }
}