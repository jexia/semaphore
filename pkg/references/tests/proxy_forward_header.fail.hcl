proxy "echo" {
  forward "caller" {
    header {
      Authorization = "{{ input:unknown }}"
    }
  }
}
