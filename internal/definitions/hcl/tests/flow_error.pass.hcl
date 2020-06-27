flow "" {
	error "proto.Error" {
		message = "{{ error:message }}"
		status = "{{ error:status }}"

        message "nested" {
            message "nested" {}
            repeated "" "" {}
        }

        repeated "" "" {
            message "nested" {}
            repeated "" "" {}
        }
	}

	on_error {
        schema = ""
		status = 401
		message = "flow error message"

        params {
            prop = ""
        }
	}
}