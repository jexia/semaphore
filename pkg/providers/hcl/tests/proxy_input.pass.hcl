proxy "echo" {
  input {}

  forward "" {}
}

proxy "ping" {
  input {
    header = ["Authorization"]
  }

  forward "" {}
}

proxy "ping" {
  input {
    options {
      key = "value"
    }
  }

  forward "" {}
}

proxy "ping" {
  input {
    params = "com.semaphore.Message"
  }

  forward "" {}
}
