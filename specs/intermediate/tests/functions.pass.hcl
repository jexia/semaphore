flow "echo" {
    input {
        message = "<string>"
    }

    output {
        message = "{{ add(input:message) }}"
    }
}