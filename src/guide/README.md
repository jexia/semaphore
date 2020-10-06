---
meta:
  - name: description
    content: Introduction
  - name: keywords
    content: semaphore getting started introduction
---

# Introduction

Semaphore is a tool to orchestrate your micro-service architecture. Requests could be manipulated passed branched to different services to be returned as a single output.
You could define request flows on top of your currently existing schema definitions.
Please check out the examples directory for more examples.

::: tip
In many of the available examples are protobuffers used. Semaphore currently supports protobuffers more official schema definitions such as Avro and XML will be added in the future
:::

```hcl
endpoint "GetUser" "http" {
    endpoint = "/user/:id"
    method = "GET"
}

flow "GetUser" {
    input "proto.Query" {}
    
    resource "user" {
        request "proto.Users" "Get" {
            id = "{{ input:id }}"
        }
    }
    
    output "proto.User" {
        name = "{{ user:name }}"
    }
}
```