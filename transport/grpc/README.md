# gRPC

Provides gRPC client and server implementations.

```hcl
endpoint "mock" "grpc" {
    package = "maestro.greeter"
	service = "Say"
    method = "Hello"
}
```
