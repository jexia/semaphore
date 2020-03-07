service "service" "http" "json" {
	host = "service.local"
	schema = "mock"
}

flow "complete" {
	input "complete" {}

	call "first" {
		request "service" "complete" {
			message = "{{ input:message }}"

			message "nested" {
				value = "{{ input:nested.value }}"
			}

			repeated "repeating" "input:repeating" {
				value = "{{input:repeating.value}}"
			}
		}
	}
}