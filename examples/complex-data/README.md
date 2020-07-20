# Complex data structures

This example shows how complex data structures could be implemented inside Semaphore.

## Getting started

To run this example you need to have the Semaphore daemon installed on your machine.

```bash
$ semaphore daemon
```

You could execute the `ComplexDataStructure` flow by executing a `POST` request on port `8080`.

```bash
$ curl 127.0.0.1:8080 -d '{
    "items": [
        {
            "id": 1,
            "name": "milk",
            "labels": ["discount", "breakfast"]
        },
        {
            "id": 2,
            "name": "bread",
            "labels": ["breakfast"]
        }
    ],
    "shipping": {
        "time": 1592550251,
        "address": {
            "street": "Kalverstraat",
            "city": "Amsterdam"
        }
    }
}'
```