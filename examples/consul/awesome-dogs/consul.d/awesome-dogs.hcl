services {
  name = "awesome-dogs"
  tags = ["web", "dogs"]
  port = 3333
  checks = [
    {
      id = "api"
      name = "HTTP API on port 3333"
      http = "http://localhost:3333/check"
      method = "get"
      interval = "5s"
      timeout = "1s"
    }
  ]
}