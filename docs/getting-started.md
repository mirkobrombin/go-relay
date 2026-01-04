# Getting Started with Go Relay

Go Relay is a type-safe, distributed job queue library designed for Go. It prioritizes developer experience (Generics) and architectural flexibility (Pluggable Brokers).

## Installation

```bash
go get github.com/mirkobrombin/go-relay
```

## Quick Start

Here is how to set up a basic in-memory processing queue.

```go
package main

import (
    "context"
    "fmt"
    "github.com/mirkobrombin/go-relay/pkg/manager"
)

type EmailJob struct {
    To      string
    Subject string
}

func main() {
    // By default, it uses an in-memory broker suitable for local dev.
    r := manager.New()

    // Notice the type safety: the second argument MUST be a function accepting (ctx, EmailJob).
    manager.Register(r, "email:send", func(ctx context.Context, job EmailJob) error {
        fmt.Printf("Sending email to %s\n", job.To)
        return nil
    })

    // This blocks, or run it in a goroutine.
    go r.Start(context.Background())

    // If you pass a struct other than EmailJob, your code won't compile.
    manager.Enqueue(context.Background(), r, "email:send", EmailJob{
        To:      "user@example.com",
        Subject: "Hello World",
    })
}
```

## Next Steps

- Explore [Architecture](architecture.md) to understand Brokers and Workers.
- Learn about [Transports](transports.md) (Redis, Memory, Warp Mesh).
