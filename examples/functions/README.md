# Functions

This example uses the functions and HTTP implementations to call the API endpoints at [jsonplaceholder.typicode.com](https://jsonplaceholder.typicode.com/).

# Getting started

You could get started by executing the `main.go` file.

```bash
$ go run main.go
```

Once Maestro is up and running you could execute the `todo` flow by calling the service on port `8080`.

```bash
$ curl 127.0.0.1:8080/ -H 'Authorization: super-secret'
```