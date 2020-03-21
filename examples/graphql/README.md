# GraphQL

This example uses the GraphQL and HTTP implementations to call the API endpoints at [jsonplaceholder.typicode.com](https://jsonplaceholder.typicode.com/).

# Getting started

You could get started by  the Maestro CLI.

```bash
$ maestro run -c config.yaml
```

Once Maestro is up and running you could execute the `todo` flow by calling the service on port `8080`.

```bash
$ curl 127.0.0.1:8080/ -d '{"query": "{latest{id}todo(id:\"2\"){title}}"}'
```