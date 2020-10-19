flow "echo" {
  input "object" {}

  resource "get" {
    request "getter" "Get" {
      message "nested" {
        name = "{{ input:nested.name }}"

        message "sub" {
          message = "hello world"
        }
      }
    }
  }

  output "object" {
    message "nested" {
      name = "<string>"
    }
  }
}
