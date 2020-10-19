service "com.semaphore" "caller" {
  transport = "http"
  codec     = "json"
  host      = ""
}

flow "echo" {
  input "com.input" {}

  resource "opening" {
    request "caller" "Open" {
      message = "{{ input:message }}"
    }
  }

  output "com.output" {
    message = "{{ input:message }}"
  }
}
