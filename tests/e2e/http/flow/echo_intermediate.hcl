endpoint "typetest" "http" {
  endpoint = "/json"
  method   = "POST"
  codec    = "json"
}

flow "typetest" {
  input "semaphore.typetest.Request" {}

  resource "echo" {
    request "semaphore.typetest.External" "Post" {}
  }

  output "semaphore.typetest.Data" {
    enum = "{{ echo:enum }}"
  }
}

// http {
//     address = ":8080"
// }

// services {
//     select "semaphore.typetest.*" {
//         host = "http://127.0.0.1:8081/"
//     }
// }