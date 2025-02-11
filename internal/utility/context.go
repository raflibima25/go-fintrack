package utility

import (
	"context"
	"time"
)

func ContextWithTimeout(ctx context.Context, timeout time.Duration) context.Context {
	ctx, _ = context.WithTimeout(ctx, timeout)
	return ctx
}
