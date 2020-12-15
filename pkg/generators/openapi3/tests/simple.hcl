protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "post"
}

flow "CreateUser" {
    input {
        payload "com.semaphore.User" {}
    }

    output {
        payload "com.semaphore.User" {}
    }
}