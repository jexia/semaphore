service "mock" "http" "json" {
	host = "mock.local"
	schema = "proto.mock"
}

flow "simple" {
	schema = "proto.mock.simple"

	call "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	schema = "proto.mock.nested"

	call "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	schema = "proto.mock.repeated"

	call "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}