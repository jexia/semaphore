flow "simple" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "simple" {
            message = "{{ input:message }}"
        }
    }
}

flow "nested" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "nested" {
            message "nested" {
                value = "{{ input:nested.value }}"
            }
        }
    }
}

flow "repeated" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "repeated" {
            repeated "repeating" "input:repeating" {
                value = "{{ input:repeating.value }}"
            }
        }
    }
}

flow "repeated_values" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "repeated_values" {
            repeated_values = "{{ input:repeating_values }}"
        }
    }
}

flow "enum" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "enum" {
            enum = "{{ input:enum }}"
        }
    }
}

flow "repeating_enum" {
    input {
        payload "com.complete.input" {}
    }

    resource "first" {
        request "mock" "repeating_enum" {
            enum = "{{ input:repeating_enum }}"
        }
    }
}
