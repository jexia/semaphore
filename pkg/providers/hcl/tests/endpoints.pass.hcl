endpoint "echo" "http" {
  random = "value"
}

endpoint "ping" "http" {
  message "random" {
    value = "message"
  }
}
