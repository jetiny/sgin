package common

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

type TokenHandler = func(c *gin.Context, token string) *UserToken
type AppModelHandler = func(c *gin.Context, appCode string) *AppModel
type LogHandler func(c *gin.Context, tableId int, opId LogOpType, primaryId int64, data any)

type BootContext struct {
	Xorm            *xorm.Engine
	Redis           *redis.Client
	SessionHandle   gin.HandlerFunc
	Logger          *logrus.Logger
	TokenHandle     TokenHandler
	AppModdelHandle AppModelHandler
	LogHandle       LogHandler
	Routes          []*Route
	Tasks           map[string][]gin.HandlerFunc
	Addr            string
	data            sync.Map
}

var gBootContext *BootContext

func NewBootContext() *BootContext {
	return &BootContext{
		Tasks:  make(map[string][]gin.HandlerFunc),
		Routes: make([]*Route, 0),
	}
}

func SetBootContext(bootContext *BootContext) {
	gBootContext = bootContext
}

func GetBootContext() *BootContext {
	return gBootContext
}

func (s *BootContext) WithRoutes(routes ...[]*Route) {
	for _, v := range routes {
		s.Routes = append(s.Routes, v...)
	}
}

func (s *BootContext) Get(key string) any {
	v, ok := s.data.Load(key)
	if ok {
		return v
	}
	return nil
}

func (s *BootContext) Load(key string) (any, bool) {
	v, ok := s.data.Load(key)
	return v, ok
}

func (s *BootContext) Set(key string, value any) {
	s.data.Store(key, value)
}

func (s *BootContext) WithTask(tasks map[string][]gin.HandlerFunc) {
	for k, v := range tasks {
		s.Tasks[k] = v
	}
}

func (s BootContext) PrintAddr() {
	Logger.Println("http://" + s.Addr)
}

type BootPlugin func(ctx *BootContext, r *gin.Engine)

func (s *BootContext) Install(r *gin.Engine, plugins ...BootPlugin) {
	for _, v := range plugins {
		v(s, r)
	}
}
