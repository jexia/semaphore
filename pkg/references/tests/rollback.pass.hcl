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

    rollback "caller" "Open" {
      message = "{{ opening:message }}"
    }
  }

  output "com.output" {
    message = "{{ opening:message }}"
  }
}
