# Maestro

Maestro is a tool to orchistrate requests inside your microservice architecture.
Maestro provides a powerfull toolset for manipulating, forwarding and returning properties from and to multiple services.

> 🚧 This project is still under construction and may be changed or updated without notice

## Getting started

All data streams inside Maestro are called flows.
A flow is able to manipulate, deconstruct and forwarded data in between calls and services.
Flows are exposed through endpoints. Flows are generic and could handle different protocols and codecs within a single flow.
All flows are strict typed through schema definitions. These schemas define the contracts provided and accepted by services.

```hcl
endpoint "checkout" "http" "json" {
    method = "POST"
    endpoint = "/checkout"
}

flow "checkout" {
    input {
        id = "<string>"
    }

    call "shipping" {
        request "warehouse" "Send" {
            user = "{{ input:id }}"
        }
    }

    output {
        status = "{{ shipping:status }}"
    }
}

service "warehouse" "grpc" "proto" {
    host = "warehouse.local"
    schema = "proto.Warehouse"
}
```