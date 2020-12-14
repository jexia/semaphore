flow "echo" {
  input "object" {}

  resource "get" {
    request "getter" "Get" {
      array = [
        "{{ input:id }}",
        "{{ input:name }}",
        "static",
      ]
    }
  }

  output "object" {
    object = {
      "message": "hello world",
      "meta": {
        "id": "{{ getter:output }}"
      }
    }
  }
}
