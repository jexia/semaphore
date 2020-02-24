caller "http" {
	base = "/v1"
}

service "test" "http" {
	host = "test.local"
	schema = "proto.test"
    codec = "proto"
}

flow "complete" {
	schema = "proto.test.complete"

	call "first" "test.complete" {
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