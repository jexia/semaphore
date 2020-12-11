flow "echo" {
  input "object" {}

  output {
    payload "object" {
      message = "{{ input:message }}"
    }
  }
}
