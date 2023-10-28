package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jetiny/sgin"
	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/uses"
)

func h(ctx *gin.Context) {
	sgin.GetCtx(ctx).Success("Ok!")
}

func p(c *gin.Context) {
	fmt.Println(isFunctionEqual(c.Handler(), p))
	sgin.GetCtx(c).Pass()
}

func isFunctionEqual(v any, i any) bool {
	return fmt.Sprintf("%v", v) == fmt.Sprintf("%v", i)
}

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
	ctx.SetTokenHandler("", func(c *gin.Context, token string) *common.UserToken {
		return &common.UserToken{
			AccessToken:      token,
			RefreshToken:     token,
			AccessExpiredAt:  time.Now().Add(time.Hour * 100),
			RefreshExpiredAt: time.Now().Add(time.Hour * 100),
		}
	})
	ctx.WithRoutes([]*common.Route{
		{Method: http.MethodPost, Path: "/user/profile/info", EnsureAuth: true, Handle: h},
		{Method: http.MethodGet, Path: "/path/*value", EnsureAuth: false, Handle: h},
	})
	ctx.Install(r, func(ctx *common.BootContext, r *gin.Engine) {
		r.GET("/pass", p)
	})
	uses.Setup(ctx, r)
	ctx.PrintAddr()
	r.Run(ctx.Addr)
}
