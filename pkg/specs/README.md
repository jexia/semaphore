# Maestro
Maestro is a tool to orchestrate your microservices by providing a powerful toolset for manipulating, forwarding and returning properties from and to multiple services.

Maestro is built on top of schema definitions and flows.
Messages are strictly typed and are type-checked. Payloads such as protobuf and JSON could be generated from the same definitions.

## Table of contents

- [Specification](#specification)
  * [Resources](#resources)
    + [Input](#input)
    + [Call](#call)
  * [Template reference](#template-reference)
  * [Message](#message)
  * [Repeated message](#repeated-message)
  * [Flow](#flow)
    + [Input](#input-1)
    + [Output](#output)
  * [Resource](#resource)
    + [Depends on](#depends-on)
    + [Options](#options)
    + [Header](#header)
    + [Request](#request)
    + [Rollback](#rollback)
  * [Proxy](#proxy)
  * [Service](#service)
    + [Options](#options)
  * [Endpoint](#endpoint)
- [Functions](#functions)

## Specification

### Resources
Objects (ex: input, call or responses) holding or representing data are called resources. Resources are only populated after they have been called.

Some resources hold multiple resource properties. The default resource property is used when no resource property is given.
The following resource properties are available:

#### Input
- **request - *default***
- header
#### Call
- request
- **response - *default***
- header

### Template reference
Templates could reference properties inside other resources. Templates are defined following the mustache template system. Templates start with the resource definition. The default resource property is used when no resource property is given.

Paths reference a property within the resource. Paths could target nested messages or repeated values.

```
{{ call.request:address.street }}
```


### Message
A message holds properties, nested messages and/or repeated messages. All of these properties could be referenced. Messages reference a schema message.
Properties
Properties hold constant values and/or references. Properties are strictly typed and use the referenced schema message for type checks. Properties could also hold references which should match the given property type.
Nested messages
You can define and use message types inside other message types, as in the following example.

```hcl
message "address" {
    message "country" {

    }
}
```
### Repeated message
Repeated messages accept two labels the first one is its alias and the second one is the resource reference. If a repeated message is kept empty the whole message is attempted to be copied.

```hcl
repeated "address" "input:address" {
    message "country" {

    }
}
```

### Flow
A flow defines a set of calls that should be called chronologically and produces an output message. Calls could reference other resources when constructing messages. All references are strictly typed. Properties are fetched from the given schema or inputs.

All flows should contain a unique name. Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
A dependency reference structure is generated within the flow which allows Maestro to figure out which calls could be called parallel to improve performance.

An optional schema could be defined which defines the request/response messages.

```hcl
flow "Logger" {
    input "schema.Object" {
    }

    resource "log" {
        request "logger" "Log" {
            message = "{{ input:message }}"
        }
    }

    output "schema.Object" {
        status = "{{ log:status }}"
        code = "{{ log:code }}"
    }   
}
```

#### Schema
A schema definition defines the input and output message types. When a flow schema is defined are the input properties (except header) ignored.

```hcl
flow "Logger" {
    schema = "exposed.Logger.Log"
}
```

#### Input
The input represents a schema definition. The schema definition defines the message format. Additional options or headers could be defined here as well.

```hcl
input "schema.Object" {
    header = ["Authorization"]
}
```

#### Output
The output acts as a message. The output could contain nested messages and repeated messages. The output could also define the response header.

```hcl
output "schema.Object" {
  header {
    Cookie = "mnomnom"
  }

  status = "{{ log:status }}"
}
```

### Resource
A call calls the given service and method. Calls could be executed synchronously or asynchronously. All calls are referencing a service method, the service should match the alias defined inside the service. The request and response schema messages are used for type definitions.
A call could contain the request headers, request body and rollback.

```hcl
# Calling service alias logger.Log
resource "log" {
  request "logger" "Log" {
    message = "{{ input:message }}"
  }
}
```

#### Depends on
Marks resources as dependencies. This could be usefull if a resource does not have a direct reference dependency.

```hcl
resource "warehouse" {
  request "warehouse" "Ship" {
    product = "{{ input:product }}"
  }
}

resource "log" {
  depends_on = ["warehouse"]

  request "logger" "Log" {
    message = "{{ input:message }}"
  }
}
```

#### Options
Options could be consumed by implementations. The defined key/values are implementation-specific.

```hcl
resource "log" {
    request "logger" "Log" {
        options {
            // HTTP method
            method = "GET"
        }
    }
}
```

#### Header
Headers are represented as keys and values. Both keys and values are strings. Values could reference properties from other resources.

```hcl
input "schema.Input" {
    header = ["Authorization"]
}

resource "log" {
    request "logger" "Log" {
        header {
            Authorization = "{{ input.header:Authorization }}"
        }
    }
}
```

#### Request
The request acts as a message. The request could contain nested messages and repeated messages.

```hcl
resource "log" {
    request "logger" "Log" {
        key = "value"
    }
}
```

#### Rollback
Rollbacks are called in a reversed chronological order when a call inside the flow fails.
All rollbacks are called async and errors are not handled.
Rollbacks consist of a call endpoint and a request message.
Rollback templates could only reference properties from any previous calls and the input.

```hcl
resource "log" {
    rollback "logger" "Log" {
        header {
            Claim = "{{ input:Claim }}"
        }
        
        message = "Something went wrong"
    }
}
```

### Proxy
A proxy streams the incoming request to the given service.
Proxies could define calls that are executed before the request body is forwarded.
The `input.request` resource is unavailable in proxy calls.
A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.

```hcl
proxy "upload" {
    resource "auth" {
        request "authenticate" "Authenticate" {
            header {
                Authorization = "{{ input.header:Authorization }}"
            }
        }
    }

    resource "logger" {
        request "logger" "Log" {
            message = "{{ auth:claim }}"
        }
    }

    forward "uploader" "File" {
        header {
            StorageKey = "{{ auth:key }}"
        }
    }
}
```

### Service
Services represent external service which could be called inside the flows.
The service name is an alias that could be referenced inside calls.
The host of the service and schema service method should be defined for each service.
The request and response message defined inside the schema are used for type definitions.
The FQN (fully qualified name) of the schema method should be used.
Each service references a caller implementation to be used.

Codec is the message format used for request and response messages.

```hcl
service "logger" {
    transport = "http"
    codec = "proto"
    host = "https://service.prod.svc.cluster.local"
}
```

#### Options
Options could be consumed by implementations. The defined key/values are implementation-specific.

```hcl
options {
    port = 8080
}
```

### Endpoint
An endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller. The name of the endpoint represents the flow which should be executed.

All servers should define their own request/response message formats.

```hcl
endpoint "users" "http" "json" {
    method = "GET"
    endpoint = "/users/:project"
    status = "202"
}
```

#### Options
Options could be consumed by implementations. The defined key/values are implementation-specific.

```hcl
options {
    port = 8080
}
```

## Functions

Custom defined functions could be configured and passed to Maestro. Functions could be called inside templates and could accept arguments and return a property as a response.
Functions could be used to preform computation on properties during runtime. Functions have read access to the entire reference store but could only write to their own stack.
A unique resource is created for each function call where all references stored during runtime are located. This resource is created during compile time and references made to the given function are automatically adjusted.

A function should always return a property where all paths are absolute. This way it is easier for other properties to reference a resource.

```
function(...<arguments>)
```

```hcl
resource "auth" {
    request "authenticate" "Authenticate" {
        header {
            Authorization = "{{ jwt(input.header:Authorization) }}"
        }
    }
}
```