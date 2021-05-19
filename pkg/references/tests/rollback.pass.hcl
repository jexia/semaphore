service "com.semaphore" "caller" {
  transport = "http"
  codec     = "json"
  host      = ""
}

flow "echo" {
    input {
        payload "com.input" {}
    }

    resource "opening" {
        request "caller" "Open" {
            message = "{{ input:message }}"
        }

        rollback "caller" "Open" {
            message = "{{ opening:message }}"
        }
    }

    output {
        payload "com.output" {
            message = "{{ opening:message }}"
        }
    }
}
