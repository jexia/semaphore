---
id: flows.requests
title: Service Requests
sidebar_label: Service Requests
slug: /flows/requests
---

Requests could be executed inside a resource or rollback block.
Requests interact with a external service through the configured transport protocol (ex: `HTTP`, `gRPC`).

The service method request and response object are used to define the request schemas.
Responses produced by the given service could be referenced inside other resources of the given flow.

```hcl
flow "CreateUser" {
    input "com.org.User" {}

    resource "user" {
        request "com.org.Users" "Create" {
          first_name = "{{ input:first_name }}"
          last_name = "{{ input:last_name }}"
          age = "{{ input:age }}"
        }

        rollback {
            request "com.org.Users" "Delete" {
              id = "{{ user:id }}"
            }
        }
    }

    output "com.org.User" {
        id = "{{ user:id }}"
        ref = "{{ input:ref }}"
    }
}
```