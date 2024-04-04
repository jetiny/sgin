package uses

import (
	"net"

	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/utils"
	"github.com/oschwald/geoip2-golang"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-xorm/xorm"
)

const (
	gHttpSuccessMessage = "Ok"
	gHttpSuccessCode    = "200"
)

type Ctx struct {
	c        *gin.Context
	pageInfo *utils.PageInfo
}

func (ctx *Ctx) PageInfo() *utils.PageInfo {
	if ctx.pageInfo != nil {
		return ctx.pageInfo
	}
	pageInfo := &utils.PageInfo{}
	err := ctx.c.ShouldBindBodyWith(pageInfo, binding.JSON)
	if err != nil {
		pageInfo.Current = 1
	}
	if pageInfo.Current == 0 {
		pageInfo.Current = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}
	ctx.pageInfo = pageInfo
	return ctx.pageInfo
}

func (ctx *Ctx) Pass() {
	ctx.Success(nil)
}

func (ctx *Ctx) Success(data any) {
	res := utils.Data[any]{
		Code:     gHttpSuccessCode,
		Data:     data,
		Message:  gHttpSuccessMessage,
		PageInfo: ctx.pageInfo,
	}
	ctx.c.JSON(200, res)
}

func (ctx *Ctx) Page(data any, total int64) {
	ctx.PageInfo().Total = total
	res := utils.Data[any]{
		Code:     gHttpSuccessCode,
		Data:     data,
		Message:  gHttpSuccessMessage,
		PageInfo: ctx.pageInfo,
	}
	ctx.c.JSON(200, res)
}

// mixed
func (ctx Ctx) Session() sessions.Session {
	return getSession(ctx.c)
}

func (ctx Ctx) Orm() *xorm.Engine {
	return getOrm(ctx.c)
}

func (ctx *Ctx) Redis() *rediCache {
	return getRedis(ctx.c)
}

func (ctx *Ctx) DataLog() *dbLog {
	return getDbLog(ctx.c)
}

func (ctx *Ctx) App() *common.AppModel {
	return getAppModel(ctx.c)
}

func (ctx *Ctx) Token() *common.UserToken {
	return getUserToken(ctx.c)
}

func (ctx *Ctx) GetToken() *common.UserToken {
	return hasUserToken(ctx.c)
}

func (ctx *Ctx) Route() *common.Route {
	return getRoute(ctx.c)
}

func (ctx *Ctx) Ip() string {
	return ctx.c.ClientIP()
}

func (ctx *Ctx) GetIpInfo() (*geoip2.City, error) {
	return getIp(ctx.c).City(net.ParseIP(ctx.c.ClientIP()))
}

func (ctx *Ctx) Stack(value ...any) {
	getErrorStack(ctx.c).Push(value...)
}

const ctxKey = "ctxKey"

func withCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Ctx{
			c: c,
		}
		c.Set(ctxKey, ctx)
		c.Next()
	}
}

func GetCtx(c *gin.Context) *Ctx {
	iter, exists := c.Get(ctxKey)
	if !exists {
		ctx := &Ctx{
			c: c,
		}
		c.Set(ctxKey, ctx)
		return ctx
	}
	return iter.(*Ctx)
}
