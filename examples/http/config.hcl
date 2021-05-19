log_level    = "$LOG_LEVEL"
protobuffers = ["./proto/*.proto"]

include = ["flow.hcl"]

http {
    address = ":8080"
}

services {
    select "com.semaphore.*" {
        host = "https://jsonplaceholder.typicode.com/"
    }
}
