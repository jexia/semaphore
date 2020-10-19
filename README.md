# Semaphore [![CI](https://github.com/jexia/semaphore/workflows/Go%20CI/badge.svg)](https://github.com/jexia/semaphore/actions?query=workflow%3A%22Go+CI%22)

----

- Chat: [Discord](https://chat.jexia.com)
- Documentation: [https://jexia.gitbook.io/semaphore/](https://jexia.gitbook.io/semaphore/)
- Go package documentation: [GoDev](https://pkg.go.dev/github.com/jexia/semaphore)

> Take control of your data, connect with anything, and expose it anywhere through protocols such as HTTP, GraphQL, and gRPC!

Create advanced and high performing data flows and expose them through endpoints over multiple protocols such as HTTP, GraphQL, and gRPC.
Create custom extensions or use the availability of custom functions and protocol implementations.

**Key features of Semaphore are:**

* **🔗 Connect with anything** Use the right tool for the job.
  Semaphore supports various protocols out of the box with
  the ability to supporting additional protocols through
  modules. Endpoints could be created to expose a single
  flow through multiple protocols.

* **🚀 Blazing fast** Semaphore scales up to your needs. Branches
  are created to execute resources concurrently. Branches are
  based on dependencies between resources made through
  references or hard coded values. Creating high-performance
  flows is almost boringly easy.

* **✅ Transactional flows** Make sure that your data stays consistent.
  Rollback data when an unexpected response is returned from
  one of your services. References to returned values could be
  made allowing to ensure that your customers have the best experience
  possible.

* **⛩️ Conditional logic** Only call services when needed.
  Conditional expressions ensure that resources are only
  executed when needed. Conditions grow to your needs.
  Whether you want to keep things simple or need to achieve
  complex goals.

* **🌍 Adapts to your environment** Semaphore integrates
  with your existing system(s). Define flows through
  simple and strict typed definitions. Use your already
  existing schema definitions such as Protobuffers. Or
  extend Semaphore with custom modules and proprietary
  software. Integrate services through flow definitions
  and create a great experience for your customers and
  your teams.
  
----

[![asciicast](https://asciinema.org/a/344280.svg)](https://asciinema.org/a/344280)

## Enterprise

Want to take your systems to the next level?
Semaphore Enterprise allows users to fully embrace the power their data flows.
Additional modules and tooling allows users to build more complex environments and helps running Semaphore in production.

Feel free to request for more information or a demo by sending us a email at:
support@jexia.com

## Documentation and Getting Started

Documentation is available at [GitBook](https://jexia.gitbook.io/semaphore/).

If you are new to Semaphore and want to get started with building flows, please
check out the available [🚀 Examples](https://github.com/jexia/semaphore/tree/master/examples).
Feel free to reach out to the community on [Discord](https://chat.jexia.com) or by opening a new issue.

Data streams inside Semaphore are defined as flows. A flow could manipulate,
deconstruct, and forwarded data in between resources. Flows are exposed through
endpoints. Flows are generic and could handle different protocols and codecs
from a single flow. All flows are strictly typed through schema definitions.
These schemas define the contracts provided and accepted by services.

Currently, are only protobuffers supported but more schema definitions are
planned to be supported in the future. Feel free to open a new issue to discuss
which schema definition you require.

```hcl
endpoint "checkout" "http" {
	endpoint = "/cart/checkout"
	method = "POST"
}

endpoint "checkout" "grpc" {
	package = "webshop.cart"
	service = "Payment"
	method = "Checkout"
}

flow "checkout" {
	input "services.Order" {}

	resource "product" {
		request "services.Warehouse" "GetProduct" {
			product = "{{ input:product }}"
		}
	}

	resource "shipping" {
		request "services.Warehouse" "Send" {
			user = "{{ input:user }}"
		}
	}

	output "services.OrderResult" {
		status = "{{ shipping:status }}"
		product = "{{ product:. }}"
	}
}
```

### Installing Semaphore

There are variouse sources available to download and install the [⚡ Semaphore CLI](https://github.com/jexia/semaphore/tree/master/cmd/semaphore). For more information and install methods please check out the [installing section](https://github.com/jexia/semaphore/tree/master/cmd/semaphore#installing).

```sh
$ curl https://raw.githubusercontent.com/jexia/semaphore/master/install.sh | sh
```

![Install Semaphore](https://user-images.githubusercontent.com/3440116/88109404-bf256800-cbaa-11ea-964e-55b089e57cd7.gif)

## Developing Semaphore

If you wish to work on Semaphore itself or any of its built-in systems, you'll
first need [Go](https://www.golang.org) installed on your machine. Go version
1.13.7+ is *required*.

For local dev first make sure Go is properly installed, including setting up a
[GOPATH](https://golang.org/doc/code.html#GOPATH). Ensure that `$GOPATH/bin` is in
your path as some distributions bundle old version of build tools. Next, clone this
repository. Semaphore uses [Go Modules](https://github.com/golang/go/wiki/Modules),
so it is recommended that you clone the repository ***outside*** of the GOPATH.
You can then download any required build tools by bootstrapping your environment:

```sh
$ make bootstrap
...
```

To compile a development version of Semaphore, run `make` or `make dev`. This will
put the Semaphore binary in the `bin` folders:

```sh
$ make dev
...
$ bin/semaphore
...
```

To run tests, type `make test`. If
this exits with exit status 0, then everything is working!

```sh
$ make test
...
```

## Contributing

Thank you for your interest in contributing to Semaphore! ❤
Check out the open projects and/or issues and feel free to join any ongoing discussion.

Everyone is welcome to contribute, whether it's in the form of code, documentation, bug reports, feature requests, or anything else. We encourage you to experiment with the project and make contributions to help evolve it to meet your needs!

See the [contributing guide](https://github.com/jexia/semaphore/blob/master/CONTRIBUTING.md) for more details.
