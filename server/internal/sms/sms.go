package sms

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"
)

// Sender 短信发送器。MVP 默认打日志，生产换阿里云/腾讯云实现。
type Sender interface {
	Send(phone, code string) error
}

type LogSender struct{}

func (LogSender) Send(phone, code string) error {
	log.Printf("[sms] to %s: code=%s", phone, code)
	return nil
}

// CodeStore 验证码存储（MVP 内存版，生产换 Redis）
type CodeStore struct {
	mu sync.Mutex
	m  map[string]codeEntry
}

type codeEntry struct {
	code string
	exp  time.Time
}

func NewCodeStore() *CodeStore { return &CodeStore{m: map[string]codeEntry{}} }

func (s *CodeStore) Set(phone, code string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[phone] = codeEntry{code: code, exp: time.Now().Add(ttl)}
}

// Verify 一次性校验：成功则删除
func (s *CodeStore) Verify(phone, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[phone]
	if !ok || time.Now().After(e.exp) || e.code != code {
		return false
	}
	delete(s.m, phone)
	return true
}

// NewCode 6 位数字验证码
func NewCode() string {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "123456"
	}
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1000000)
}
