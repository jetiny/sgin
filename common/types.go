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
	TokenKey         string    `json:"tokenKey"`
	IsJwt            bool      `json:"isJwt"`
	userData         any
}

func (s *UserToken) SetUserData(data any) {
	s.userData = data
}

func (s *UserToken) GetUserData() any {
	return s.userData
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
	return GetAppConfig(s, env.KeyName(), env.Value().(T))
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

type RouteHandler func(ctx *gin.Context)

func (s RouteHandler) Route() *Route {
	return &Route{
		Handle: s,
	}
}

// 路由定义
type Route struct {
	Name       string
	Label      string
	Path       string
	Method     string
	EnsureAuth bool
	NoAppCode  bool
	Handle     RouteHandler
	TokenKey   string
}

func (s *Route) WithPath(value string) *Route {
	s.Path = value
	return s
}

func (s *Route) WithMethod(value string) *Route {
	s.Method = value
	return s
}

func (s *Route) WithName(value string) *Route {
	s.Name = value
	return s
}

func (s *Route) WithLabel(value string) *Route {
	s.Label = value
	return s
}

func (s *Route) WithTokenKey(value string) *Route {
	s.TokenKey = value
	return s
}

func (s *Route) WithAuth(value bool) *Route {
	s.EnsureAuth = value
	return s
}

func (s *Route) WithNoAppCode(value bool) *Route {
	s.NoAppCode = value
	return s
}

type RouteQS struct {
	route      Route
	Name       *string
	Label      *string
	Method     *string
	EnsureAuth *bool
	NoAppCode  *bool
	TokenKey   *string
}

type RouteBuilder struct {
	RouteQS
	routes []*Route
}

func NewRouteBuidler() *RouteBuilder {
	s := &RouteBuilder{
		routes: make([]*Route, 0),
	}
	return s
}

func (s *RouteBuilder) AddRoute(path string, handler RouteHandler) *RouteBuilder {
	s.routes = append(s.routes, handler.Route().WithPath(path))
	return s
}

func (s *RouteBuilder) WithName(value string) *RouteBuilder {
	s.route.Name = value
	s.Name = &s.route.Name
	return s
}

func (s *RouteBuilder) WithLabel(value string) *RouteBuilder {
	s.route.Label = value
	s.Label = &s.route.Label
	return s
}

func (s *RouteBuilder) WithMethod(value string) *RouteBuilder {
	s.route.Method = value
	s.Method = &s.route.Method
	return s
}

func (s *RouteBuilder) WithAuth(value bool) *RouteBuilder {
	s.route.EnsureAuth = value
	s.EnsureAuth = &s.route.EnsureAuth
	return s
}

func (s *RouteBuilder) WithNoAppCode(value bool) *RouteBuilder {
	s.route.NoAppCode = value
	s.NoAppCode = &s.route.NoAppCode
	return s
}

func (s *RouteBuilder) WithTokenKey(value string) *RouteBuilder {
	s.route.TokenKey = value
	s.TokenKey = &s.route.TokenKey
	return s
}

func (s *RouteBuilder) Result() ([]*Route, *RouteQS) {
	return s.routes, &s.RouteQS
}
