flow "echo" {
  resource "unknown" {
    request "caller" "Open" {
      message = "{{ input.header:Unkown }}"
    }
  }
}
