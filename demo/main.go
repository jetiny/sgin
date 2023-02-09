package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jetiny/sgin"
	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/uses"
)

func main() {
	ctx := sgin.MustBootAll()
	// ctx := sgin.MustBoot(base.BootWithNone)
	r := gin.Default()
	ctx.AppModdelHandle = func(c *gin.Context, appCode string) *common.AppModel {
		return &common.AppModel{
			AppId:   1,
			AppCode: "8888",
			AppName: "8888",
		}
	}
	ctx.TokenHandle = func(c *gin.Context, token string) *common.UserToken {
		return &common.UserToken{
			AccessToken:      token,
			RefreshToken:     token,
			AccessExpiredAt:  time.Now().Add(time.Hour * 100),
			RefreshExpiredAt: time.Now().Add(time.Hour * 100),
		}
	}
	ctx.WithRoutes([]*common.Route{
		{Method: http.MethodPost, Path: "/user/profile/info", EnsureAuth: true, Handle: func(ctx *gin.Context) {
			sgin.GetCtx(ctx).Pass()
		}},
	})
	ctx.Install(r, func(ctx *common.BootContext, r *gin.Engine) {
		r.GET("/pass", func(c *gin.Context) {
			sgin.GetCtx(c).Pass()
		})
	})
	uses.Setup(ctx, r)
	ctx.PrintAddr()
	r.Run(ctx.Addr)
}
