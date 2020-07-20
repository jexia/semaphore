# Micro

The micro transport implements a wrapper around `go-micro` services.
The configured transporter and registry are used inside the flow manager.

```go
// gRPC service constructor
service := grpc.NewService()

client, err := semaphore.New(
        semaphore.WithCodec(json.NewConstructor()),
        semaphore.WithCaller(micro.New("micro-grpc", service)),
)
```