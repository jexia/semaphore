# Maestro

Maestro is a tool to orchistrate requests inside your microservice architecture.
Maestro provides a powerfull toolset for manipulating, forwarding and returning properties from and to multiple services.

> 🚧 This project is still under construction and may be changed or updated without notice

# Getting started

All exposed endpoints are defined as flows.
A flow could manipulate, deconstruct and pass data in between calls and services.
All message and services are defined inside a schema format (currently only protobuf is supported).
Flows are exposed through endpoints. Each endpoint could contain server specific configurations.

```hcl
flow "checkout" {
	input {
        id = "<string>"
        customer = "<string>"

        message "address" {
            city = "<string>"
        }
	}

	call "prepare" "warehouse.Prepare" {
		request {
			cart = "{{ input:id }}"

            rollback "warehouse.Cancel" {
                cart = "{{ input:id }}"
            }
		}
	}

	call "send" "shipping.Send" {
		request {
			order = "{{ prepare:order }}"
			customer = "{{ input:id }}"
			city = "{{ input:address.city }}"
		}
	}

    output {
        ref = "{{ prepare:order }}"
        status = "{{ send:url }}"
    }
}
```