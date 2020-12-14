proxy "echo" {
  forward "uploader" {}
}

proxy "ping" {
  forward "uploader" {
    header {
      cookie = "mnomnom"
    }
  }
}
