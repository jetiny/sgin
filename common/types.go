package common

import (
	"encoding/json"
	"time"

	"github.com/jetiny/sgin/utils"

	"github.com/gin-gonic/gin"
)

type UserToken struct {
	UserInfo         UserInfo  `json:"userInfo"`
	AccessToken      string    `json:"accessToken"`
	RefreshToken     string    `json:"refreshToken"`
	CreatedAt        time.Time `json:"createdAt"`
	AccessExpiredAt  time.Time `json:"accessExpiredAt"`
	RefreshExpiredAt time.Time `json:"refreshExpiredAt"`
}

func (s UserToken) IsTokenExpired() bool {
	return s.AccessExpiredAt.Unix() < time.Now().Unix()
}

type UserInfo struct {
	UserId   int64  `json:"userId,string,omitempty" `
	Nickname string `json:"nickname,omitempty" `
	Avatar   string `json:"avatar,omitempty" `
}

type AppModel struct {
	AppId   int64  `json:"appId,string"` // 应用id
	AppCode string `json:"appCode"`      // 应用标识
	AppName string `json:"appName"`      // 应用名称
	config  map[string]any
}

func (s AppModel) Config() map[string]any {
	return s.config
}

func (s *AppModel) SetConfig(config string) {
	if config != "" {
		json.Unmarshal([]byte(config), &s.config)
	} else {
		s.config = map[string]any{}
	}
}

func (s AppModel) Get(key string) (any, bool) {
	v, b := s.config[key]
	return v, b
}

func GetEnvConfig[T int | uint | int64 | uint64 | bool | string | []any | map[string]any](s *AppModel, env *utils.EnvGetter) T {
	value, ok := s.Get(env.KeyName())
	if ok {
		return value.(T)
	}
	return env.DefaultValue().(T)
}

func GetAppConfig[T int | uint | int64 | uint64 | bool | string | []any | map[string]any](s *AppModel, key string, defaultValue T) T {
	value, ok := s.Get(key)
	if ok {
		return value.(T)
	}
	return defaultValue
}

// Dblog
type LogOpType int16

const (
	LogOpNormal LogOpType = 0 // 常规
	LogOpInsert LogOpType = 1 // 更新记录
	LogOpUpdate LogOpType = 2 // 更新记录
	LogOpDelete LogOpType = 3 // 删除记录
	LogOpDrop   LogOpType = 4 // 彻底删除记录
)

// 路由定义
type Route struct {
	Name       string
	Path       string
	Method     string
	EnsureAuth bool
	Handle     gin.HandlerFunc
}
