flow "simple" {
	input {
		payload "proto.Simple" {}
	}

	resource "first" {
		request "proto.mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input {
		payload "proto.Message" {}
	}

	resource "first" {
		request "proto.mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input {
		payload "proto.Message" {}
	}

	resource "first" {
		request "proto.mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}

flow "repeated_values" {
	input {
		payload "proto.Message" {}
	}

	resource "first" {
		request "proto.mock" "repeated_values" {
			repeating = "{{ input:repeating_values }}"
		}
	}
}