proxy "echo" {
  input {
    params = "com.input"
  }

  resource "mock" {
    rollback "caller" "Open" {
      message = "{{ input:message }}"
    }
  }

  forward "caller" {}
}
