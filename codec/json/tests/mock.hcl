service "mock" "http" {
	host = "mock.local"
	schema = "mock"
    codec = "json"
}

flow "simple" {
	schema = "mock.simple"

	call "first" "mock.simple" {
		request {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	schema = "mock.nested"

	call "first" "mock.nested" {
		request {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	schema = "mock.repeated"

	call "first" "mock.repeated" {
		request {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}