package common

import (
	"context"
	"time"
)

func NewContext(duration *time.Duration) (context.Context, context.CancelFunc) {
	if duration == nil {
		var timeout = time.Second * 10
		duration = &timeout
	}
	return context.WithTimeout(context.Background(), *duration)
}
