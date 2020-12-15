flow "" {
  resource "" {
    error {
      // TODO: allow references
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
      status  = 401
      message = "node error message"
    }
  }
}
