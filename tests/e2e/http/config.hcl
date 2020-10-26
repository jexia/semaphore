log_level = "$LOG_LEVEL"

protobuffers = ["./proto/*.proto"]

include = ["./flow/echo.hcl"]

grpc {
  address = ":50051"
}

http {
  address = ":8080"
}
