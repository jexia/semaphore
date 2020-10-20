---
id: flows.conditions
title: Conditions
sidebar_label: Conditions
slug: /flows/conditions
---

Conditional logic allows for resources only to be executed if a given condition returns `true`.
Complex expressions could be created and property references could be used.

```hcl
flow "IsAdmin" {
	input "org.Input" {}

	if "{{ input:is_admin }}" {
		resource "query" {
			request "proto.Service" "ThrowError" {}
		}
	}

	output "org.Output" {}
}
```

## Complex conditions

Complex conditions could be defined with multiple operators and using multiple resources.

```hcl
if "({{ input:has_property }} && {{ input:age }} > 18) && {{ input:is_admin }}" {
	resource "query" {
		request "proto.Service" "ThrowError" {}
	}
}
```