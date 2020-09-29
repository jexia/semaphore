proxy "echo" {
  input {
    params = "unexpected"
  }

  forward "caller" {}
}
