<h1 align="left"><a href="https://jexia.com"><img src="https://user-images.githubusercontent.com/3440116/77702983-019eb580-6fba-11ea-8d2c-f6a6b8e60cbd.jpg" alt="Jexia"></a> Maestro</h1>

<p align="left">
  <a href="https://pkg.go.dev/github.com/jexia/maestro"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="GoDev"></a>
  <a href="https://github.com/jexia/maestro/actions?query=workflow%3A%22Go+CI%22"><img src="https://github.com/jexia/maestro/workflows/Go%20CI/badge.svg" alt="Go CI"></a>
  <a href="https://goreportcard.com/report/github.com/jexia/maestro"><img src="https://goreportcard.com/badge/github.com/jexia/maestro" alt="Go report"></a>
  <a href="https://jexia.gitbook.io/maestro/"><img src="https://img.shields.io/badge/docs-gitbook-green" alt="Gitbook"></a>
  <a href="https://discord.gg/q54Q8GH"><img src="https://img.shields.io/badge/chat-on%20discord-7289da.svg?sanitize=true" alt="Chat on Discord"></a>
</p>

Maestro is a feature-rich service orchistrator. Create advanced data flows and expose them through endpoints.
Have full controll over your exposed endpoints, expose single flows for multiple protocols such as gRPC and GraphQL.
Maestro adapts to your environment, create custom extentions or use the available of custom functions and protocol implementations.

## Table of contents

1. [Using Maestro](#using-maestro)
1. [Getting started](#getting-started)
1. [Contributing](#contributing)

## Using Maestro

Maestro could be used in a wide variety of cases. It could be used to let teams have full controll over their exposed endpoints.
Create SAGA patterns to autimatically rollback requests on failure. Allow users to implement your product with their tools of choice.
We are excited to see how you will implement Maestro in your architecture.

* **Gateway**: Maestro redefines the gateway. Expose a single flow through the multiple protocols without changing any of your services.

* **Scalable**: You are able to scale Maestro up to your needs. All calls within a flow are executed in the most optimal path possible. Branches are created to execute calls [concurrently](https://github.com/jexia/maestro/tree/master/flow) from one another when possible.

* **SAGA patterns**: Define rollbacks inside your flows in the case of failure. Rollbacks are automatically executed if a request fails. Rollbacks could reference data received from other services.

* **E2E testing**: Expose your internal e2e tests through any protocol. Deploy a Maestro instance to expose internal endpoints without exposing them to the public.

## Getting started

1. [⚡ Daemon](https://github.com/jexia/maestro/tree/master/cmd/daemon)
1. [🚀 Examples](https://github.com/jexia/maestro/tree/master/examples)
1. [📚 Documentation](https://jexia.gitbook.io/maestro/)

You could download the daemon from source or most commonly used package managers. Or pull one of the available [docker images](https://hub.docker.com/r/jexiacom/maestro).

```bash
docker pull jexiacom/maestro
```

---

Data streams inside Maestro are defined inside flows.
A flow could manipulate, deconstruct and forwarded data in between calls and services.
Flows are exposed through endpoints. Flows are generic and could handle different protocols and codecs within a single flow.
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
