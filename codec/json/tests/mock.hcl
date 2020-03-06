service "mock" "http" "json" {
	host = "mock.local"
	schema = "mock"
}

flow "simple" {
	input "simple" {}

	call "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "nested" {}

	call "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "repeated" {}

	call "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}