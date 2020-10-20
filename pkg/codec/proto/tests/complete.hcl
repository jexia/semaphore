flow "complete" {
  input "proto.Message" {}

  resource "first" {
    request "proto.test" "complete" {
      message = "{{ input:message }}"

      message "nested" {
        value = "{{ input:nested.value }}"
      }

      repeated "repeating" "input:repeating" {
        value = "{{input:repeating.value}}"
      }

      repeating_values = "{{ input:repeating_values }}"
      status           = "{{ input:status }}"
      repeating_status = "{{ input:repeating_status }}"
    }
  }
}
