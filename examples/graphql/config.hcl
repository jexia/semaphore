log_level = "$LOG_LEVEL"
protobuffers = ["./annotations", "./proto/*.proto"]

include = ["flow.hcl"]

graphql {
    address = ":8080"
}

prometheus {
    address = ":5050"
}

services {
    select "proto.*" {
        host = "https://jsonplaceholder.typicode.com/"
    }
}
