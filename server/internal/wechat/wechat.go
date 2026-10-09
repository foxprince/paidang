package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Session struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 只在服务端调用，secret 永不下发到小程序
func Code2Session(appID, secret, code string) (*Session, error) {
	u := "https://api.weixin.qq.com/sns/jscode2session?" + url.Values{
		"appid":      {appID},
		"secret":     {secret},
		"js_code":    {code},
		"grant_type": {"authorization_code"},
	}.Encode()
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var s Session
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	if s.ErrCode != 0 {
		return nil, fmt.Errorf("wechat: %d %s", s.ErrCode, s.ErrMsg)
	}
	return &s, nil
}
