flow "echo" {
    input {
        message "nested" {
            name = "<string>"
        }
    }

    call "get" "getter.Get" {
        request {
            message "nested" {
                name = "{{ input:nested.name }}"

                message "sub" {
                    message = "hello world"
                }
            }
        }
    }
}