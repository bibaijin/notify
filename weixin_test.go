package notify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWeixinNotify(t *testing.T) {
	cases := []struct {
		corpID    string
		appID     int
		appSecret string
		users     []string
		message   string
	}{
		{
			// 这里需要换成真实信息
			"fake",
			0,
			"fake",
			[]string{
				"fake",
			},
			"Hello, fake.",
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		w := NewWeixin(c.corpID, c.appID, c.appSecret, zap.L())
		err := w.Notify(ctx, c.users, c.message, zap.L())
		assert.Nil(t, err)
	}
}
