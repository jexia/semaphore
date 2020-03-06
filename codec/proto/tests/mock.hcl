service "mock" "http" "json" {
	host = "mock.local"
	schema = "proto.mock"
}

flow "simple" {
	input "proto.Simple" {}

	call "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "proto.Message" {}

	call "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "proto.Message" {}

	call "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}