# gRPC

This example uses the gRPC implementation to allow semaphore to call a simple greeting service.

## Getting started

To run this example you need to have Go 1.13> and the Semaphore CLI installed on your machine.
First start the service by simply executing the files inside the service directory.

```bash
$ cd service
$ go run .
```

Start Semaphore to expose the greeting service on port `8080`.

```bash
$ cd gateway
$ semaphore daemon -f config.hcl
```

You could execute the `greeter` flow by using [evans](https://github.com/ktr0731/evans) or by calling the HTTP endpoint.

```bash
$ evans -r -p 50051
$ curl 127.0.0.1:8080 -d '{"name":"world"}'
```