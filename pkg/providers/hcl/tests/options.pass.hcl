log_level = "debug"

protobuffers = ["$PROTO_IMPORT"]

grpc {
  address = "$GRPC"
}

http {
  address = "$HTTP"
}

graphql {
  address = "$GRAPHQL"
}
