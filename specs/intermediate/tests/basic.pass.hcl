flow "echo" {
    input {
        message = "<string>"
    }

    output {
        message = "{{ input:message }}"
    }
}