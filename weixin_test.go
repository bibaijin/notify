package notify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeixin(t *testing.T) {
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
		w := NewWeixin(c.corpID, c.appID, c.appSecret)
		token, expiresIn, err := w.GetToken(ctx)
		require.Nil(t, err)
		assert.NotZero(t, token)
		assert.NotZero(t, expiresIn)

		err = w.Notify(ctx, token, c.users, c.message)
		require.Nil(t, err)
	}
}
