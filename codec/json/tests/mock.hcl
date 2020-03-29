flow "simple" {
	input "complete" {}

	resource "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "complete" {}

	resource "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "complete" {}

	resource "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}

flow "repeated_values" {
	input "complete" {}

	resource "first" {
		request "mock" "repeated_values" {
			repeated_values = "{{ input:repeating_values }}"
		}
	}
}