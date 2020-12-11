flow "" {
  error {
    // TODO: fixme
    status = 400 // "{{ error:status }}"

    payload "proto.Error" {
      message = "{{ error:message }}"

      message "nested" {
        message  "nested" {}
        repeated "" "" {}
      }

      repeated "" "" {
        message  "nested" {}
        repeated "" "" {}
      }
    }
  }

  on_error {
    schema  = "com.Schema"
    status  = 401
    message = "flow error message"

    params {
      prop = ""
    }
  }
}
