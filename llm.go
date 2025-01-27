package main

import (
	"context"
	"fmt"
	"time"
)


// llm generates tokens (e.g., "Token1", "Token2") at a specified interval.
// It respects the provided context for timeouts or cancellations.
func llm(ctx context.Context, interval time.Duration, totalTokens int) <-chan string {
	tokenStream := make(chan string)

	go func() {
		defer close(tokenStream) // Ensure the channel is closed when done

		for i := 1; i <= totalTokens; i++ {
			select {
			case <-ctx.Done(): // Stop generating if context is canceled
				return
			case tokenStream <- fmt.Sprintf("Token-%d", i): // Send the token
				time.Sleep(interval) // Wait for the specified interval
			}
		}
	}()

	return tokenStream
}
