package notify

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockNotifier 是 Notifier 的一个 mock
type MockNotifier struct {
	mock.Mock
}

// NewMock 返回初始化后的 mock
func NewMock() *MockNotifier {
	return &MockNotifier{}
}

// Notify 发送通知
func (m *MockNotifier) Notify(ctx context.Context, users []string, message string) error {
	args := m.Called()
	return args.Error(0)
}
