package notify

import (
	"context"
	"time"
)

// Notifier 表示通知接口
type Notifier interface {
	GetToken(ctx context.Context) (token string, expiresIn time.Duration, err error)
	Notify(ctx context.Context, token string, users []string, message string) error
}
