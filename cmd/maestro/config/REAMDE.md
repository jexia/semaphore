# Config

CLI configurations could be stored inside a yaml configuration file.
A configuration file could be referenced when executing a CLI command.

```yaml
level: "debug"
http:
    address: ":8080"
graphql:
    address: ":9090"
grpc:
    address: ":50051"
protobuffers:
- "../annotations"
- "./*.proto"
flows:
- "./*.hcl"
```