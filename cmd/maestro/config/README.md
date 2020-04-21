# Config

CLI configurations could be stored inside a `HCL` configuration file.
A configuration file could be referenced when executing a CLI command.
Environment variables could be used inside the `HCL` definitions.

```hcl
log_level = "$LOG_LEVEL"
protobuffers = ["../../annotations", "./proto/*.proto"]

include = ["flow.hcl"]

graphql {
    address = "$GRAPHQL_ADDRESS"
}

http {
    address = "$HTTP_ADDRESS"
}

services {
    select "proto.users.*" {
        host = "api.jexia.com"
    }

    select "proto.projects.*" {
        host = "api.jexia.com"
    }
}
```