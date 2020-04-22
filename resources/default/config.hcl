log_level = "info"
protobuffers = ["./proto/**.proto"]

include = ["./flows/**.hcl"]

http {
    address = ":80"
}
