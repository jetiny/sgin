package common

import (
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
}

var gBootContext *BootContext

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

func (s *BootContext) WithTask(tasks map[string][]gin.HandlerFunc) {
	for k, v := range tasks {
		s.Tasks[k] = v
	}
}

type BootPlugin func(ctx *BootContext, r *gin.Engine)

func (s *BootContext) Install(r *gin.Engine, plugins ...BootPlugin) {
	for _, v := range plugins {
		v(s, r)
	}
}
