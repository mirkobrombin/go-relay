package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mirkobrombin/go-relay/v2/pkg/manager"
)

type Notification struct {
	UserID  int
	Message string
}

func main() {
	r := manager.New()

	manager.Register(r, "notify", func(ctx context.Context, n Notification) error {
		fmt.Printf("[Worker] Sending notification to User %d: %s\n", n.UserID, n.Message)
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Println("Relay started. Waiting for jobs...")
		if err := r.Start(ctx); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Enqueueing jobs...")
	_ = manager.Enqueue(ctx, r, "notify", Notification{UserID: 101, Message: "Hello World"})
	_ = manager.Enqueue(ctx, r, "notify", Notification{UserID: 102, Message: "Order Shipped"})
	_ = manager.Enqueue(ctx, r, "notify", Notification{UserID: 103, Message: "System Alert"})

	time.Sleep(1 * time.Second)
	fmt.Println("Done.")
}
