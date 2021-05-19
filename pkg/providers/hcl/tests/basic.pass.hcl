flow "echo" {
  input {
    payload "object" {}
  }

  output {
    payload "object" {
      message = "{{ input:message }}"
    }
  }
}
