package uses

import (
	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

func Setup(ctx *common.BootContext, r *gin.Engine) {
	r.Use(withCors())
	r.Use(withCtx())
	if ctx.SessionHandle != nil {
		r.Use(ctx.SessionHandle)
	}
	if ctx.Logger != nil {
		r.Use(withHttpLog(ctx.Logger))
	}
	r.Use(withRecovery())
	r.Use(withGinError())
	if ctx.Xorm != nil {
		r.Use(withOrm(ctx.Xorm))
	}
	if ctx.Redis != nil {
		r.Use(withRedis(ctx.Redis))
	}
	if ctx.LogHandle != nil {
		r.Use(withDbLog(ctx.LogHandle))
	}
	if ctx.AppModelHandle != nil {
		r.Use(withAppModel(ctx.AppModelHandle))
	}
	if ctx.Tasks != nil {
		r.Use(withTask(ctx.Tasks))
	}
	routes := ctx.Routes
	if routes == nil {
		routes = make([]*common.Route, 0)
	}
	withRoute(r, routes, ctx.TokenHandlers)
}
