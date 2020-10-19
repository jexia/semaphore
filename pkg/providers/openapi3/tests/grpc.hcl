protobuffers = ["./*.proto"]

endpoint "CreateUser" "grpc" {
  service = "users"
  method = "create"
}

flow "CreateUser" {
  input "com.semaphore.User" {}

  output "com.semaphore.User" {}
}