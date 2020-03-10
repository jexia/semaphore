flow "simple" {
	input "proto.Simple" {}

	call "first" {
		request "proto.mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "proto.Message" {}

	call "first" {
		request "proto.mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "proto.Message" {}

	call "first" {
		request "proto.mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}