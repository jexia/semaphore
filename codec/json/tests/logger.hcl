caller "http" {
	base = "/v1"
}

service "logger" "http" {
	host = "logger.local"
	schema = "Logger"
    codec = "proto"
}

flow "logger" {
	schema = "Logger.Append"

	call "logging" "logger.Append" {
		request {
			message = "{{ input:message }}"

			message "nested" {
				value = "{{ input:nested.value }}"
			}

			repeated "repeating" "input:repeating" {
				value = "{{input:repeating.value}}"
			}
		}
	}

	output {
        id = "{{ logging:id }}"
    }
}