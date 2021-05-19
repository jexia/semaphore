protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "post"
}

flow "CreateUser" {
    input {
        header = ["Authorization", "X-IP"]

        payload "com.semaphore.User" {}
    }

    output {
        payload "com.semaphore.User" {}
    }
}