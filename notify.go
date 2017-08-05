package notify

import (
	"context"

	"go.uber.org/zap"
)

// Notifier 表示通知接口
type Notifier interface {
	Notify(ctx context.Context, users []string, message string, logger *zap.Logger) error
}
