service "service" "http" {
	host = "service.local"
	schema = "test"
    codec = "json"
}

flow "complete" {
	schema = "test.complete"

	call "first" "service.complete" {
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
}