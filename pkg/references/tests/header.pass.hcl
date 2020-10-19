service "com.semaphore" "caller" {
  transport = "http"
  codec     = "json"
  host      = ""
}

flow "echo" {
  input "com.input" {
    header = ["Authorization"]
  }

  resource "opening" {
    request "caller" "Open" {
      header {
        Authorization = "{{ input.header:Authorization }}"
      }
    }
  }
}
