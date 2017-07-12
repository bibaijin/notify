package notify

import (
	"context"
)

// Notifier 表示通知接口
type Notifier interface {
	Notify(ctx context.Context, users []string, message string) error
}
