endpoint "echo" "http" "json" {
    random = "value"
}

endpoint "ping" "http" "proto" {
    message "random" {
        value = "message"
    }
}