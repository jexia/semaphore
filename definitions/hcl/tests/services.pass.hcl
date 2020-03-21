service "com.maestro" "auth" "http" "json" {
    host = "https://auth.com"
}

service "com.maestro" "users" "http" "proto" {
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