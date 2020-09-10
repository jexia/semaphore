# HTTP

This example holds a simple Semaphore HTTP implementation and shows the usage of the Semaphore Go API.

## Getting started

To run this example you need to have Go 1.13> installed on your machine.
First start the service by simply executing the files inside the service directory.

```bash
$ semaphore daemon
```

You could execute the `FetchLatestProject` flow by executing a `GET` request on port `8080`.

```bash
$ curl 127.0.0.1:8080
```