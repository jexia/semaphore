flow "simple" {
	input "proto.Simple" {}

	resource "first" {
		request "proto.mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "proto.Message" {}

	resource "first" {
		request "proto.mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "proto.Message" {}

	resource "first" {
		request "proto.mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}

flow "repeated_values" {
	input "proto.Message" {}

	resource "first" {
		request "proto.mock" "repeated_values" {
			repeating = "{{ input:repeating_values }}"
		}
	}
}

flow "oneof" {
	input "proto.Message" {}

	resource "first" {
		request "proto.mock" "one_of" {
			oneof = "{{ input:oneof }}"
		}
	}
}