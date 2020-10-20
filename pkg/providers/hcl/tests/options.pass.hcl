log_level = "debug"

protobuffers = ["$PROTO_IMPORT"]
openapi3 = ["$OPENAPI3_IMPORT"]

grpc {
  address = "$GRPC"
}

http {
  address = "$HTTP"
}

graphql {
  address = "$GRAPHQL"
}
