---
id: flows.errors
title: Error handling
sidebar_label: Error handling
slug: /flows/errors
---

Unexpected errors could happen.
Custom error messages, status codes and response objects could be defined inside your flows.
These objects could reference properties and be overriden with constant values.
Error handling consists out of two blocks.

`error` which defines the custom response object.
These objects are returned to the user if the protocol allows for dynamic error objects (such as `HTTP`).
In case when no dynamic error objects could be returned are the error message and status code used.

`on_error` allows for the definitions of parameters (params) and to override the message and status properties.
Optionally could a schema be defined. This schema is used to decode the received message.
The default error properties (message and status), error params and other resources could be referenced inside the `on_error` and error blocks.

Check out the [errors example](https://github.com/jexia/semaphore/tree/master/examples/error-handling) for a more hands on approach on how to define custom errors.

```hcl
error "com.Schema" {
    message "meta" {
        status = "{{ error:status }}"
    }

    message = "{{ error:message }}"
}

flow "greeter" {
    # copies the global error object if not set

    resource "echo" {
        # copies the flow error object if not set
    }
}
```

Error objects are cloned from the above scope when no error message has been defined.
This allows for global error messages to be defined and shared across flows.
