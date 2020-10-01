proxy "echo" {
  input {
    params = "com.input"
  }

  forward "caller" {}
}
