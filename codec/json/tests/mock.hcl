service "mock" "http" "json" {
	host = "mock.local"
	schema = "mock"
}

flow "simple" {
	schema = "mock.simple"

	call "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	schema = "mock.nested"

	call "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	schema = "mock.repeated"

	call "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}