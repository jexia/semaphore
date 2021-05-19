protobuffers = ["./*.proto"]

endpoint "CreateUser" "grpc" {
    service = "users"
    method = "create"
}

flow "CreateUser" {
    input {
        payload "com.semaphore.User" {}
    }

    output {
        payload "com.semaphore.User" {}
    }
}