package config

import "os"

// 所有密钥走环境变量，永不进仓库
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	WechatAppID  string
	WechatSecret string
	// 订阅消息模板 ID（微信公众平台申请）
	TplNewBooking    string
	TplBookingResult string
	TplBookingDone   string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	return &Config{
		Port:         getenv("PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		WechatAppID:  os.Getenv("WECHAT_APPID"),
		WechatSecret: os.Getenv("WECHAT_SECRET"),

		TplNewBooking:    os.Getenv("WECHAT_TPL_NEW_BOOKING"),
		TplBookingResult: os.Getenv("WECHAT_TPL_BOOKING_RESULT"),
		TplBookingDone:   os.Getenv("WECHAT_TPL_BOOKING_DONE"),
	}
}
