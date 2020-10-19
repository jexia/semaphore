service "com.semaphore" "caller" {
  transport = "http"
  codec     = "json"
  host      = ""
}

proxy "echo" {
  resource "opening" {
    request "caller" "Open" {}
  }

  resource "reference" {
    request "caller" "Open" {
      message = "{{ opening:message }}"
    }
  }

  forward "caller" {}
}
