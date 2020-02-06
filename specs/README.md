# Maestro
Maestro is a tool to orchestrate your microservices by providing a powerful toolset for manipulating, forwarding and returning properties from and to multiple services.

Maestro is built on top of proto buffers and flows.
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
    + [Dependency](#dependency)
  * [Call](#call-1)
    + [Header](#header)
    + [Request](#request)
    + [Rollback](#rollback)
  * [Proxy](#proxy)
  * [Service](#service)
  * [Caller](#caller)
  * [Endpoint](#endpoint)

## Specification

### Resources
Objects (ex: input, call or responses) holding or representing data are called resources. Resources are only populated after they have been called.

Some resources hold multiple resource properties. The default resource property is used when no resource property is given.
The following resource properties are available:

#### Input
- **request - *default***
- header
#### Call
- **response - *default***
- header

### Template reference
Templates could reference properties inside other resources. Templates are defined following the mustache template system. Templates start with the resource definition. The default resource property is used when no resource property is given.

Paths reference a property within the resource. Paths could target nested messages or repeated values.

```
{{ call.request:address.street }}
```


### Message
A message holds properties, nested messages and/or repeated messages. All of these properties could be referenced. Messages reference a protobuf message.
Properties
Properties hold constant values and/or references. Properties are strictly typed and use the referenced protobuf message for type checks. Properties could also hold references which should match the given property type.
Nested messages
You can define and use message types inside other message types, as in the following example.

```hcl
message "address" {
    message "country" {

    }
}
```
### Repeated message
Repeated messages are messages which are repeated. Nested messages could be defined inside repeated messages. Repeated messages accept two labels the first one is its alias and the second one is the resource reference. If a repeated message is kept empty the whole message is attempted to be copies. Repeated messages could not be defined inside a repeated message.

```hcl
repeated "address" "{{ input:address }}" {
    message "country" {

    }
}
```

### Flow
A flow defines a set of calls that should be called chronologically and produces an output message. Calls could reference other resources when constructing messages. All references are strictly typed. Properties are fetched from the given proto buffers or inputs.

All flows should contain a unique name. Calls are nested inside of flows and contain two labels, a unique name within the flow and the service and method to be called.
A dependency reference structure is generated within the flow which allows Maestro to figure out which calls could be called parallel to improve performance.

```hcl
flow "Logger" {
    input {
        message = "string"
    }

    call "log" "logger.Log" {
        request {
            message = "{{ input:message }}"
        }
    }

    output {
        status = "{{ log:status }}"
        code = "{{ log:code }}"
    }   
}
```

#### Input
The input acts as a message. The input could contain nested messages and repeated messages. Input properties could reference types and or constant values. Input types are defined by wrapping the type inside angle brackets.

```hcl
input {
    type = "sync"
    message = "<string>"
}
```
#### Output
The output acts as a message. The output could contain nested messages and repeated messages. The output could also define the response header.

```hcl
output {
  header {
    Cookie = "mnomnom"
  }

  status = "{{ log:status }}"
}
```

#### Dependency
Dependencies are flows that need to be called before the given flow is executed. Dependencies could have other dependencies which have to be called.

```hcl
flow "GetUsers" {
    dependency = ["Auth", "HasGetPolicy"]
}
```

### Call
A call calls the given service and method. Calls could be executed synchronously or asynchronously. All calls are referencing a service method, the service should match the alias defined inside the service. The request and response messages are used for type definitions.
A call could contain the request headers, request body, rollback, and the execution type.

```hcl
# Calling service alias logger.Log
call "log" "logger.Log" {
  type = "sync" # default value

  request {
    message = "{{ input:message }}"
  }
}
```

#### Header
Headers are represented as keys and values. Both keys and values are strings. Values could reference properties from other resources.

```hcl
header {
    Authorization = "{{ input.header:Authorization }}"
}
```

#### Request
The request acts as a message. The request could contain nested messages and repeated messages.

#### Rollback
Rollbacks are called in a reversed chronological order when a call inside the flow fails.
All rollbacks are called async and errors are not handled.
Rollbacks consist of a call endpoint and a request message.
Rollback templates could only reference properties from any previous calls and the input.

```hcl
rollback "logger.Log" {
    header {
        Claim = "{{ input:Claim }}"
    }
    
    request {
        message = "Something went wrong while"
    }
}
```

### Proxy
A proxy streams the incoming request to the given service.
Proxies could define calls that are executed before the request body is forwarded.
The `input.request` resource could not reference within calls since it is not parsed.
A proxy forward could ideally be used for file uploads or large messages which could not be stored in memory.

```hcl
proxy "upload" {
    call "auth" "authenticate.Authenticate" {
        header {
            Authorization = "{{ input.header:Authorization }}"
        }

        request {}
    }

    call "logger" "logger.Log" {
        request {
            message = "{{ auth:claim }}"
        }
    }

    forward "uploader" "uploader.File" {
        header {
            StorageKey = "{{ auth:key }}"
        }
    }

    output {
        id = "{{ uploader:id }}"
    }
}
```

### Service
Services define an external service which could be called inside the flows.
The service name is an alias which could be referenced inside calls.
The host of the service and proto service method should be defined for each service.
The request and response message defined inside the proto buffers are used for type definitions.
The FQN (fully qualified name) of the proto method should be used.
Each service references a caller implementation to be used.

```hcl
service "logger" "http" {
    host = "https://service.prod.svc.cluster.local"
    proto = "proto.Logger"
}
```

### Caller
Represents a caller implementation. Each implementation has to be configured and defined before running the service. All values are passed as attributes (from the hcl2 package) to the callers. These attributes could be used for configuration purposes

```hcl
caller "http" {
  base = "/v1"
}
```

### Endpoint
An endpoint exposes a flow. Endpoints are not parsed by Maestro and have custom implementations in each caller. The name of the endpoint represents the flow which should be executed.

```hcl
endpoint "users" {
    http "GET" "/users/:project" {
        status = "202"
    }
}
```
