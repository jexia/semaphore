log_level = "$LOG_LEVEL"
protobuffers = ["../../../", "../proto/*.proto"]

include = ["flow.hcl"]

http {
    address = ":9090"
}
