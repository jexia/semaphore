flow "echo" {
    input {
        message "nested" {
            name = "<string>"
        }
    }

    call "get" {
        request "getter" "Get" {
            message "nested" {
                name = "{{ input:nested.name }}"

                message "sub" {
                    message = "hello world"
                }
            }
        }
    }
}