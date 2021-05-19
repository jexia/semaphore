service "com.semaphore" "caller" {
    transport = "http"
    codec     = "json"
    host      = ""
}

flow "echo" {
    input {
        header = ["Authorization"]

        payload "com.input" {}
    }

    resource "opening" {
        request "caller" "Open" {
            header {
                Authorization = "{{ input.header:Authorization }}"
            }
        }
    }
}
