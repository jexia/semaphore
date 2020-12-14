# gRPC

Provides gRPC client and server implementations.

```hcl
endpoint "mock" "grpc" {
    package = "semaphore.greeter"
	service = "Say"
    method = "Hello"
}
```
