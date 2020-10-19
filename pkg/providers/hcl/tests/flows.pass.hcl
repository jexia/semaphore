flow "echo" {
  resource "sample" {
    request "service" "method" {
      options {
        sample = "value"
      }
    }
  }
}

flow "ping" {}
