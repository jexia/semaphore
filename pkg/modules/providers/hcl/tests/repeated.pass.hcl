flow "echo" {
  input "object" {}

  resource "get" {
    request "getter" "Get" {
      repeated "nested" "input:nested" {
        name = "{{ input:nested.name }}"

        repeated "sub" "input:nested.sub" {
          message = "hello world"
        }
      }
    }
  }

  output "object" {
    repeated "nested" "input:nested" {
      name = "<string>"

      repeated "sub" "input:nested.sub" {
        message = "<string>"
      }
    }
  }
}

proxy "echo" {
  resource "get" {
    request "getter" "Get" {
      repeated "nested" "input:nested" {
        name = "{{ input:nested.name }}"

        repeated "sub" "input:nested.sub" {
          message = "hello world"
        }
      }
    }
  }

  forward "" {}
}
