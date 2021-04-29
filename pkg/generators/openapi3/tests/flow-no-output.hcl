protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
  endpoint = "/user"
  method = "post"
}

flow "CreateUser" {
  input "com.semaphore.User" {}
}