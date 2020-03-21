# gRPC

This example uses the `go-micro` gRPC implementation to allow maestro to call a simple greeting service.

## Getting started

To run this example you need to have Go 1.13> installed on your machine.
First start the service by simply executing the files inside the service directory.

```bash
$ cd service
$ go run .
```

Start Maestro to expose the greeting service on port `8080`.

```bash
$ cd maestro
$ maestro run
```

You could execute the `greeter` flow by executing a `POST` request on port `8080`.

```bash
$ curl 127.0.0.1:8080 -d '{"name":"world"}'
```