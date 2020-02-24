service "mock" "http" {
	host = "mock.local"
	schema = "proto.mock"
    codec = "json"
}

flow "simple" {
	schema = "proto.mock.simple"

	call "first" "mock.simple" {
		request {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	schema = "proto.mock.nested"

	call "first" "mock.nested" {
		request {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	schema = "proto.mock.repeated"

	call "first" "mock.repeated" {
		request {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}