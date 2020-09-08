flow "condition" {
  resource "first" {}

  if "{{ first:id }} == {{ first:name }}" {
    if "condition" {
      resource "" {}
    }

    resources {
      sample = ""
    }
  }
}
