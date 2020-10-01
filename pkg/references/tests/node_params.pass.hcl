flow "echo" {
  input "com.input" {}

  resource "opening" {
    request "caller" "Open" {
      params {
        message = "{{ input:message }}"
      }
    }
  }
}
