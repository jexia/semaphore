# Consul

This example holds a simple Semaphore HTTP implementation with a service reslved by Consul service discovery

## Getting started

To run this example you need to have Go 1.13> installed on your machine.

Prepare the environment:

1. Run consul `consul agent -dev --config-dir=./awesome-dogs/consul.d -ui`
2. Run awesome-dogs service: `go run ./awesome-dogs/main.go`
3. Run semaphore `semaphore daemon`

Now, you can execute the `FetchPets` flow by executing a `GET` request on port `8080`.

```bash
$ curl 127.0.0.1:8080
```