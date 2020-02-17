service "auth" "http" {
    host = "https://auth.com"
    schema = "proto.Auth"
    codec = "json"
}

service "users" "http" {
    host = "https://users.com"
    schema = "proto.Users"
    codec = "proto"
}