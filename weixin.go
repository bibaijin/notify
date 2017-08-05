package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	weixinGetTokenURL    = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	weixinSendMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	expiresInMargin      = 600 * time.Second
)

// weixin 实现了向微信发送通知
type weixin struct {
	corpID    string
	appID     int
	appSecret string
	token     string
}

// NewWeixin 返回初始化后的 wexin
func NewWeixin(corpID string, appID int, appSecret string, logger *zap.Logger) Notifier {
	w := weixin{
		corpID:    corpID,
		appID:     appID,
		appSecret: appSecret,
	}

	ctx := context.Background()
	token, _, _ := w.getToken(ctx, logger)
	w.token = token
	go w.updateToken(ctx, logger)

	return &w
}

func (w weixin) Notify(ctx context.Context, users []string, message string, logger *zap.Logger) error {
	params := url.Values{}
	params.Add("access_token", w.token)
	url := fmt.Sprintf("%s?%s", weixinSendMessageURL, params.Encode())

	body := map[string]interface{}{
		"touser":  strings.Join(users, "|"),
		"msgtype": "text",
		"agentid": w.appID,
		"text": map[string]interface{}{
			"content": message,
		},
		"safe": 0,
	}
	bs, err := json.Marshal(body)
	if err != nil {
		logger.Error("json.Marshal() failed.",
			zap.Any("body", body),
			zap.Error(err),
		)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		logger.Error("http.Post() failed.",
			zap.String("URL", url),
			zap.Error(err),
		)
		return err
	}
	logger.Debug("http.Post() done.",
		zap.String("URL", url),
		zap.Any("response", resp),
	)

	var data sendMessageResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if data.ErrCode != 0 {
		return fmt.Errorf("%v", data.ErrMsg)
	}

	return nil
}

type sendMessageResponse struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

func (w weixin) getToken(ctx context.Context, logger *zap.Logger) (token string, expiresIn time.Duration, err error) {
	params := url.Values{}
	params.Add("corpid", w.corpID)
	params.Add("corpsecret", w.appSecret)

	url := fmt.Sprintf("%s?%s", weixinGetTokenURL, params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("http.Get() failed.",
			zap.String("URL", url),
			zap.Error(err),
		)
		return "", 0, err
	}
	logger.Debug("http.Get() done.",
		zap.String("URL", url),
		zap.Any("response", resp),
	)

	var data getTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", 0, err
	}

	if data.ErrCode != 0 {
		return "", 0, fmt.Errorf("%v", data.ErrMsg)
	}

	return data.AccessToken, time.Duration(data.ExpiresIn) * time.Second, nil
}

type getTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (w *weixin) updateToken(ctx context.Context, logger *zap.Logger) {
	expiresIn := expiresInMargin + 60*time.Second
	for {
		select {
		case <-ctx.Done():
			logger.Info("updateToken() cancelled.")
			return
		case <-time.After(expiresIn - expiresInMargin):
			token, e, err := w.getToken(ctx, logger)
			if err != nil {
				logger.Error("w.GetToken() failed.", zap.Error(err))
				expiresIn = expiresInMargin + 60*time.Second
				continue
			}

			expiresIn = e
			w.token = token
		}
	}
}
