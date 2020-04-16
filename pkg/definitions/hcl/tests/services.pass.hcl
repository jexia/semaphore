service "com.maestro" "auth" {
    transport = "http"
    codec = "json"
    host = "https://auth.com"
}

service "com.maestro" "auth" {
}

service "com.maestro" "users" {
    transport = "http"
    codec = "proto"
    host = "https://users.com"

    method "Add" {
        request = "proto.Request"
        response = "proto.Response"
    }

    method "Delete" {
        request = "proto.Request"
        response = "proto.Response"
    }
}