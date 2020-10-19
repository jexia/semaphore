# Error handling

This example holds a simple Semaphore error handling implementation.

## Getting started

To run this example you need to have the Semaphore daemon installed on your machine.
First start the service by simply starting the Semaphore daemon.

```bash
$ semaphore daemon
```

You could execute one of the flows by sending a `GET` request on port `8080`.

```bash
$ curl 127.0.0.1:8080/ # global error
$ curl 127.0.0.1:8080/flow # flow error
$ curl 127.0.0.1:8080/flow/node # node error
```