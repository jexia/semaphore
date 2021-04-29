protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
  endpoint = "/user"
  method = "post"
}

flow "CreateUser" {
  output "com.semaphore.User" {}
}