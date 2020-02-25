endpoint "echo" "http" {
    codec = "json"
    random = "value"
}

endpoint "ping" "http" {
    codec = "proto"

    message "random" {
        value = "message"
    }
}