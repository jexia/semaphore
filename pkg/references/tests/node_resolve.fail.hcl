flow "echo" {
    resource "unkown" {
      request "caller" "Open" {
        message = "{{ input.header:Unkown }}"
      }
    }
}