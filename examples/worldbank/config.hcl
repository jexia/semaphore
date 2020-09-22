log_level = "DEBUG"
protobuffers = ["./proto/*.proto"]

include = ["flow.hcl"]

http {
    address = ":8080"
}

services {
    select "com.semaphore.*" {
        host = "http://api.worldbank.org/"
    }
}
