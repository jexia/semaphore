flow "echo" {
    input {
        payload "com.input" {}
    }

    resource "opening" {
        request "caller" "Open" {
            params {
                message = "{{ input:message }}"
            }
        }
    }
}
