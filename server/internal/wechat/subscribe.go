package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SubscribeClient 微信订阅消息。失败只打日志，不阻断业务。
type SubscribeClient struct {
	appID  string
	secret string

	mu    sync.Mutex
	token string
	exp   time.Time
}

func NewSubscribeClient(appID, secret string) *SubscribeClient {
	return &SubscribeClient{appID: appID, secret: secret}
}

func (c *SubscribeClient) accessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.exp) {
		return c.token, nil
	}
	resp, err := http.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		c.appID, c.secret))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r struct {
		Token   string `json:"access_token"`
		Expires int    `json:"expires_in"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.ErrCode != 0 {
		return "", fmt.Errorf("wechat token: %d %s", r.ErrCode, r.ErrMsg)
	}
	c.token = r.Token
	c.exp = time.Now().Add(time.Duration(r.Expires-300) * time.Second)
	return c.token, nil
}

// Send 发送订阅消息。data 形如 {"thing1": {"value": "..."}}
func (c *SubscribeClient) Send(openid, templateID string, data map[string]map[string]string) error {
	token, err := c.accessToken()
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{
		"touser":      openid,
		"template_id": templateID,
		"data":        data,
	})
	resp, err := http.Post(
		"https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="+token,
		"application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var r struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.ErrCode != 0 {
		return fmt.Errorf("wechat subscribe: %d %s", r.ErrCode, r.ErrMsg)
	}
	return nil
}

// Str 快捷构造文本字段
func Str(v string) map[string]string { return map[string]string{"value": v} }
