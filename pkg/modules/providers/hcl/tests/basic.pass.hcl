flow "echo" {
  input "object" {}

  output "object" {
    message = "{{ input:message }}"
  }
}
