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

	"github.com/bibaijin/notify/log"
)

const (
	weixinGetTokenURL    = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	weixinSendMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// weixin 实现了向微信发送通知
type weixin struct {
	corpID    string
	appID     int
	appSecret string
	token     string
}

// NewWeixin 返回初始化后的 wexin
func NewWeixin(corpID string, appID int, appSecret string) Notifier {
	w := weixin{
		corpID:    corpID,
		appID:     appID,
		appSecret: appSecret,
	}

	ctx := context.Background()
	token, _, _ := w.getToken(ctx)
	w.token = token
	go w.updateToken(ctx)

	return &w
}

func (w weixin) Notify(ctx context.Context, users []string, message string) error {
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
		log.Errorf(ctx, "json.Marshal(%+v) failed, error: %v.", body, err)
		return err
	}

	log.Infof(ctx, "http.Post(%v)....", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		log.Errorf(ctx, "http.Post(%v) failed, error: %v.", url, err)
	}

	var data sendMessageResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	log.Infof(ctx, "http.Post(%v) done, response: %+v.", url, data)

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

func (w weixin) getToken(ctx context.Context) (token string, expiresIn time.Duration, err error) {
	params := url.Values{}
	params.Add("corpid", w.corpID)
	params.Add("corpsecret", w.appSecret)

	url := fmt.Sprintf("%s?%s", weixinGetTokenURL, params.Encode())
	log.Infof(ctx, "http.Get(%v)....", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf(ctx, "http.Get(%v) failed, error: %v.", url, err)
		return "", 0, err
	}

	var data getTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", 0, err
	}

	log.Infof(ctx, "http.Get(%v) done, response: %+v.", url, data)

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

func (w weixin) updateToken(ctx context.Context) {
	var token string
	expiresIn := 101 * time.Second
	var err error
	for {
		select {
		case <-ctx.Done():
			log.Infof(ctx, "updateToken() cancelled.")
			return
		case <-time.After(expiresIn - 100*time.Second):
			token, expiresIn, err = w.getToken(ctx)
			if err != nil {
				log.Errorf(ctx, "w.GetToken() failed, error: %v.", err)
				continue
			}

			w.token = token
		}
	}
}
