service "auth" "http" "json" {
    host = "https://auth.com"
    schema = "proto.Auth"
}

service "users" "http" "proto" {
    host = "https://users.com"
    schema = "proto.Users"
}