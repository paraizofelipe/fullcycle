package ctxlog

import (
	"context"
	"errors"
	"log"
)

func LogDeadline(ctx context.Context, err error, label string) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Printf("%s timeout: %v", label, err)
	}
}
