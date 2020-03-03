# Maestro [![GoDoc](https://godoc.org/github.com/jexia/maestro?status.svg)](https://godoc.org/github.com/jexia/maestro) ![Go CI](https://github.com/jexia/maestro/workflows/Go%20CI/badge.svg)

Maestro is a tool to orchestrate requests inside your microservice architecture.
A request could be manipulated passed branched to different services to be returned as a single output.

The key features of Maestro are:

* **Call branching**: All calls within a flow are executed [concurrently](https://github.com/jexia/maestro/tree/master/flow) from one another. Dependencies between calls are created through references or when specified.

* **SAGA patterns**: Rollbacks are easily implemented and automatically executed when an unexpected error is thrown during execution. Rollbacks could reference data received from other services.

## Getting started

1. [🚀 Examples](https://github.com/jexia/maestro/tree/master/examples)
2. [📚 Documentation](https://godoc.org/github.com/jexia/maestro)

---

All data streams inside Maestro are called flows.
A flow can manipulate, deconstruct and forwarded data in between calls and services.
Flows are exposed through endpoints. Flows are generic and could handle different protocols and codecs within a single flow.
All flows are strictly typed through schema definitions. These schemas define the contracts provided and accepted by services.

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

## Contributing

Thank you for your interest in contributing to Maestro! ❤
Check out the open projects and/or issues and feel free to join any ongoing discussion.

Everyone is welcome to contribute, whether it's in the form of code, documentation, bug reports, feature requests, or anything else. We encourage you to experiment with the project and make contributions to help evolve it to meet your needs!

See the [contributing guide](https://github.com/jexia/maestro/blob/master/CONTRIBUTING.md) for more details.
