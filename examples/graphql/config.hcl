log_level = "$LOG_LEVEL"
protobuffers = ["../../annotations", "./proto/*.proto"]

include = ["flow.hcl", "endpoints.$ENV.hcl", "services.$ENV.hcl"]

graphql {
    address = "$GRAPHQL_ADDRESS"
}

http {
    address = "$HTTP_ADDRESS"
}

services {
    select "proto.kerberos.*" {
        host = "api.jexia.com"
    }

    select "proto.andvari.*" {
        host = "api.jexia.com"
    }
}
