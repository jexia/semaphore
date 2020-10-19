service "com.semaphore" "caller" {
  transport = "http"
  codec     = "json"
  host      = ""
}

flow "echo" {
  input "com.input" {}

  if "{{ input:message }}" {
    resource "opening" {
      request "caller" "Open" {
        message = "{{ input:message }}"
      }
    }
  }

  output "com.output" {
    message = "{{ opening:message }}"
  }
}
