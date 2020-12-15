protobuffers = ["./*.proto"]

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "get"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "post"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "put"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "patch"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "delete"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "options"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "head"
}

endpoint "CreateUser" "http" {
    endpoint = "/user"
    method = "trace"
}

flow "CreateUser" {
    input {
        payload "com.semaphore.User" {}
    }

    output {
        payload "com.semaphore.User" {}
    }
}