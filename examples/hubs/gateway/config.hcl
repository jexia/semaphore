log_level = "$LOG_LEVEL"
protobuffers = ["../../../annotations", "../proto/*.proto"]

include = ["flow.hcl"]

http {
    address = ":8080"
}
