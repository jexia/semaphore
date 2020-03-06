service "test" "http" "proto" {
	host = "test.local"
	schema = "proto.test"
}

flow "complete" {
	input "proto.Message" {}

	call "first" {
		request "test" "complete" {
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