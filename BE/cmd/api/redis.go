package main

import (
	"context"
	"fmt"
	"time"
)

func (app *application) storeActivationToken(key string, userID string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.redis.Set(ctx, key, userID, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store activation token: %w", err)
	}
	return nil
}
