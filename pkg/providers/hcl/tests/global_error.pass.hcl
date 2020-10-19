error "proto.Error" {
  message = "{{ error:message }}"
  status  = "{{ error:status }}"

  message "nested" {
    message  "nested" {}
    repeated ""       ""       {}
  }

  repeated "" "" {
    message  "nested" {}
    repeated ""       ""       {}
  }
}
