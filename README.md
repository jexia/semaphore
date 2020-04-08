<h1 align="center">Maestro <a href="https://jexia.com"><img src="https://user-images.githubusercontent.com/3440116/77702983-019eb580-6fba-11ea-8d2c-f6a6b8e60cbd.jpg" alt="Jexia"></a></h1>

<p align="center">
  <a href="https://pkg.go.dev/github.com/jexia/maestro"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="GoDev"></a>
  <a href="https://github.com/jexia/maestro/actions?query=workflow%3A%22Go+CI%22"><img src="https://github.com/jexia/maestro/workflows/Go%20CI/badge.svg" alt="Go CI"></a>
  <a href="https://goreportcard.com/report/github.com/jexia/maestro"><img src="https://goreportcard.com/badge/github.com/jexia/maestro" alt="Go report"></a>
  <a href="https://jexia.gitbook.io/maestro/"><img src="https://img.shields.io/badge/docs-gitbook-green" alt="Gitbook"></a>
  <a href="https://discord.gg/q54Q8GH"><img src="https://img.shields.io/badge/chat-on%20discord-7289da.svg?sanitize=true" alt="Chat on Discord"></a>
</p>

<img align="center" src="https://user-images.githubusercontent.com/3440116/77703025-154a1c00-6fba-11ea-9515-71156bcda177.png">

Maestro is a tool to orchestrate requests inside your microservice architecture.
A request could be manipulated passed branched to different services to be returned as a single output.

The key features of Maestro are:

* **Call branching**: All calls within a flow are executed [concurrently](https://github.com/jexia/maestro/tree/master/flow) from one another. Dependencies between calls are created through references or when specified.

* **SAGA patterns**: Rollbacks are easily implemented and automatically executed when an unexpected error is thrown during execution. Rollbacks could reference data received from other services.

* **Proxy forwarding**: Allows support for streaming protocols such as websockets and to run Maestro instances in front of each other. Allow your own team to be in charge of their own flow definitions. Check out the [hubs example](https://github.com/jexia/maestro/tree/master/examples/hubs) for more information.

## Getting started

1. [⚡ CLI](https://github.com/jexia/maestro/tree/master/cmd/maestro)
1. [🚀 Examples](https://github.com/jexia/maestro/tree/master/examples)
1. [📚 Documentation](https://jexia.gitbook.io/maestro/)

You could download the CLI from source or most commonly used package managers. Or pull one of the available docker images.

```bash
docker pull docker.pkg.github.com/jexia/maestro/cli:latest
```

---

All data streams inside Maestro are called flows.
A flow can manipulate, deconstruct and forwarded data in between calls and services.
Flows are exposed through endpoints. Flows are generic and could handle different transports and codecs within a single flow.
All flows are strictly typed through schema definitions. These schemas define the contracts provided and accepted by services.

```hcl
endpoint "checkout" "http" {
    method = "POST"
    endpoint = "/checkout"
    codec = "json"
}

endpoint "checkout" "graphql" {
    path = "payment"
    base = "mutation"
}

endpoint "checkout" "grpc" {
    package = "webshop.cart"
    service = "Payment"
    method = "Checkout"
}

flow "checkout" {
    input "schema.Object" {
    }

    resource "shipping" {
        request "package.warehouse" "Send" {
            user = "{{ input:id }}"
        }
    }

    output "schema.Object" {
        status = "{{ shipping:status }}"
    }
}
```

## Contributing

Thank you for your interest in contributing to Maestro! ❤
Check out the open projects and/or issues and feel free to join any ongoing discussion.

Everyone is welcome to contribute, whether it's in the form of code, documentation, bug reports, feature requests, or anything else. We encourage you to experiment with the project and make contributions to help evolve it to meet your needs!

See the [contributing guide](https://github.com/jexia/maestro/blob/master/CONTRIBUTING.md) for more details.
