flow "simple" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "simple" {
      message = "{{ input:message }}"
    }
  }
}

flow "nested" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "nested" {
      message "nested" {
        value = "{{ input:nested.value }}"
      }
    }
  }
}

flow "repeated" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "repeated" {
      repeated "repeating" "input:repeating" {
        value = "{{ input:repeating.value }}"
      }
    }
  }
}

flow "repeated_values" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "repeated_values" {
      repeated_values = "{{ input:repeating_values }}"
    }
  }
}

flow "enum" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "enum" {
      enum = "{{ input:enum }}"
    }
  }
}

flow "repeating_enum" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "repeating_enum" {
      enum = "{{ input:repeating_enum }}"
    }
  }
}

flow "one_of" {
  input "com.complete.input" {}

  resource "first" {
    request "mock" "one_of" {
      oneof = "{{ input:oneof }}"
    }
  }
}