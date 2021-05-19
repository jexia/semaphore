service "com.semaphore" "caller" {
    transport = "http"
    codec     = "json"
    host      = ""
}

flow "echo" {
    input {
        payload "com.input" {}
    }

    if "{{ input:message }}" {
        resource "opening" {
            request "caller" "Open" {
                message = "{{ input:message }}"
            }
        }
    }

    output {
        payload "com.output" {
            message = "{{ opening:message }}"
        }
    }
}
