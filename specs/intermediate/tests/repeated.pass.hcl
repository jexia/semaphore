flow "echo" {
    input {
        repeated "nested" {
            name = "<string>"

            repeated "sub" {
                message = "<string>"
            }
        }
    }

    call "get" {
        request "getter" "Get" {
            repeated "nested" "{{ input:nested }}" {
                name = "{{ input:nested.name }}"

                repeated "sub" "{{ input:nested.sub }}" {
                    message = "hello world"
                }
            }
        }
    }
}