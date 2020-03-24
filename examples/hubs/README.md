# Hubs

This example shows how multiple Maestro hubs could be used to forward requests.
In this example a gateway is forwarding all requests from `/v1` to a hub which holds exposed endpoints.
Before the request is forwarded to the hub is a simple (placeholder) check preformed to validate the request.
If one of the checks fails is the proxy forward not executed.

## Getting started

To run this example you need to have the Maestro CLI installed on your machine.
First start the gateway on port `8080`.

```bash
$ cd gateway
$ maestro run -c config.yaml
```

Start the Maestro hub to expose the users service on port `9090`.

```bash
$ cd hub
$ maestro run -c config.yaml
```

You could execute the `user` flow by executing a `GET` request on port `8080/v1/user`.

```bash
$ curl 127.0.0.1:8080/v1/user'
```