---
id: flows.rollbacks
title: Rollbacks
sidebar_label: Rollbacks
slug: /flows/rollbacks
---

Rollbacks are called in a reversed chronological order when a call inside the flow fails. All rollbacks are called async and any unexpected errors are ignored. Rollback templates could only reference properties from any previous executed resource, the error resource, and input.

```hcl
resource "log" {
    rollback "logger" "Log" {
        header {
            trace = "{{ trace:id }}"
        }
        
        message = "{{ error:message }}"
    }
}
```