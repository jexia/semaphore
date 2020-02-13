service "auth" "http" {
    host = "https://auth.com"
    schema = "proto.Auth"
}

service "users" "http" {
    host = "https://users.com"
    schema = "proto.Users"
}