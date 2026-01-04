# Go Relay

**Go Relay** is a distributed, type-safe background job processing library for Go.

It acts as a bridge between your application logic and your distributed infrastructure (Warp Mesh, Redis, etc.), providing a clean, generic-based API for handling async tasks.

## Key Features

- **Type-Safe**: No more `interface{}` casting. Payload types are checked at compile time.
- **Backend Agnostic**: Switch from In-Memory to Redis or P2P Mesh without changing your code.
- **Zero-Config**: Starts with a sensible In-Memory default for immediate productivity.

## Documentation

- **[Getting Started](docs/getting-started.md)**: Your first Job in 30 seconds.
- **[Architecture](docs/architecture.md)**: How Relay works under the hood.
- **[Transports](docs/transports.md)**: Configure Memory, Redis, or Warp Mesh backends.

## Quick Example

```go
type Notification struct {
    ReqID string
}

// Register
relay.Register(r, "notify", func(ctx context.Context, n Notification) error {
    return sendPush(n.ReqID)
})

// Enqueue
relay.Enqueue(ctx, r, "notify", Notification{ReqID: "123"})
```

## License

MIT License.
