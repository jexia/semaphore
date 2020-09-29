flow "complete" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "complete" {
      message = "{{ input:message }}"

      message "nested" {
        value = "{{ input:nested.value }}"
      }

      repeated "repeating" "input:repeating" {
        value = "{{input:repeating.value}}"
      }

      repeating_values = "{{ input:repeating_values }}"
      enum             = "{{ input:enum }}"
      repeating_enum   = "{{ input:repeating_enum }}"
    }
  }
}
