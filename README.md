<h1 align="left"><a href="https://jexia.com"><img src="https://user-images.githubusercontent.com/3440116/77702983-019eb580-6fba-11ea-8d2c-f6a6b8e60cbd.jpg" alt="Jexia"></a> Semaphore</h1>

<p align="left">
  <a href="https://pkg.go.dev/github.com/jexia/semaphore"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="GoDev"></a>
  <a href="https://github.com/jexia/semaphore/actions?query=workflow%3A%22Go+CI%22"><img src="https://github.com/jexia/semaphore/workflows/Go%20CI/badge.svg" alt="Go CI"></a>
  <a href="https://goreportcard.com/report/github.com/jexia/semaphore"><img src="https://goreportcard.com/badge/github.com/jexia/semaphore" alt="Go report"></a>
  <a href="https://jexia.gitbook.io/semaphore/"><img src="https://img.shields.io/badge/docs-gitbook-green" alt="Gitbook"></a>
  <a href="https://chat.jexia.com"><img src="https://img.shields.io/badge/chat-on%20discord-7289da.svg?sanitize=true" alt="Chat on Discord"></a>
  <a href="https://codecov.io/gh/jexia/semaphore"><img src="https://codecov.io/gh/jexia/semaphore/branch/master/graph/badge.svg" alt="Code test coverage"></a>
</p>

Semaphore is a feature-rich service orchestrator. Create advanced data flows and expose them through endpoints.
Have full control over your exposed endpoints, expose single flows for multiple protocols such as gRPC and GraphQL.
Semaphore adapts to your environment. Create custom extensions or use the availability of custom functions and protocol implementations.

[![asciicast](https://asciinema.org/a/344280.svg)](https://asciinema.org/a/344280)

## Adapts to your environment

Semaphore integrates with your existing system(s). Define flows through simple and strict typed definitions. Use your already existing schema definitions such as Protobuffers. Or extend Semaphore with custom modules and proprietary software. Integrate services through flow definitions and create a great experience for your customers and your teams.

## Table of contents

1. [Using Semaphore](#using-semaphore)
1. [Getting started](#getting-started)
1. [Contributing](#contributing)

## Using Semaphore

Semaphore could be used in a wide variety of cases. It could be used to let teams have full control over their exposed endpoints.
Create SAGA patterns to automatically rollback requests on failure.
Allow users to implement your product with their tools of choice. We are excited to see how you will implement Semaphore in your architecture.

* **Gateway**: Semaphore redefines the gateway. Expose a single flow through multiple protocols without changing any of your services.

* **Scalable**: You can scale Semaphore up to your needs. All calls within a flow are executed in the most optimal path possible. Branches are created to execute calls [concurrently](https://github.com/jexia/semaphore/tree/master/internal/flow) from one another when possible.

* **SAGA patterns**: Define rollbacks inside your flows in the case of failure. Rollbacks are automatically executed if a request fails. Rollbacks could reference data received from other services.

* **E2E testing**: Expose your internal e2e tests through any protocol. Deploy a Semaphore instance to expose internal endpoints without exposing them to the public.

## Getting started

1. [⚡ CLI](https://github.com/jexia/semaphore/tree/master/cmd/semaphore)
1. [🚀 Examples](https://github.com/jexia/semaphore/tree/master/examples)
1. [📚 Documentation](https://jexia.gitbook.io/semaphore/)

### Install

There are variouse install methods available. You could download and install the daemon from source or most commonly used package managers. For more information and install methods please check out the [CLI](https://github.com/jexia/semaphore/tree/master/cmd/semaphore#installing).

```sh
$ curl https://raw.githubusercontent.com/jexia/semaphore/master/install.sh | sh
```

---

Data streams inside Semaphore are defined as flows.
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

Thank you for your interest in contributing to Semaphore! ❤
Check out the open projects and/or issues and feel free to join any ongoing discussion.

Everyone is welcome to contribute, whether it's in the form of code, documentation, bug reports, feature requests, or anything else. We encourage you to experiment with the project and make contributions to help evolve it to meet your needs!

See the [contributing guide](https://github.com/jexia/semaphore/blob/master/CONTRIBUTING.md) for more details.
