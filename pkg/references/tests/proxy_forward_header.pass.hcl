proxy "echo" {
  forward "caller" {
    header {
      Authorization = "Bearer"
    }
  }
}
