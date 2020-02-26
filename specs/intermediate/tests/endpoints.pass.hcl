endpoint "echo" "http" {
    codec = "json"
    
    options {
        random = "value"
    }
}

endpoint "ping" "http" {
    codec = "proto"

    options {
        message "random" {
            value = "message"
        }
    }
}