flow "simple" {
	input "complete" {}

	call "first" {
		request "mock" "simple" {
			message = "{{ input:message }}"
		}
	}
}

flow "nested" {
	input "complete" {}

	call "first" {
		request "mock" "nested" {
			message "nested" {
				value = "{{ input:nested.value }}"
			}
		}
	}
}

flow "repeated" {
	input "complete" {}

	call "first" {
		request "mock" "repeated" {
			repeated "repeating" "input:repeating" {
				value = "{{ input:repeating.value }}"
			}
		}
	}
}