---
id: functions
title: Functions
sidebar_label: Functions
slug: /functions
---

Functions could be used to preform computation on properties during runtime. Functions have read access to the entire reference store but could only write to their own stack.
A unique resource is created for each function call where all references stored during runtime are located. This resource is created during compile time and references made to the given function are automatically adjusted.

A function should always return a property where all paths are absolute. This way it is easier for other properties to reference a resource.

```
function(...<arguments>)
```

Functions could be called inside templates and could accept arguments and return a property as a response.
A collection of predefined functions is included inside the Semaphore CLI.

```hcl
resource "auth" {
    request "com.project" "Authenticate" {
        header {
            Authorization = "{{ jwt(input.header:Authorization) }}"
        }
    }
}
```