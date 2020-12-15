flow "echo" {
  output {
    payload "com.input" {
      message = "{{ input.header:Unknown }}"
    }
  }
}
