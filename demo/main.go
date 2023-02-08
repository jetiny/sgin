package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jetiny/sgin"
	"github.com/jetiny/sgin/base"
	"github.com/jetiny/sgin/common"
)

func main() {
	ctx := sgin.MustBoot(base.BootWithNone)
	r := gin.Default()
	ctx.Install(r, func(ctx *common.BootContext, r *gin.Engine) {
		r.GET("/pass", func(c *gin.Context) {
			sgin.GetCtx(c).Pass()
		})
		r.GET("/success", func(c *gin.Context) {
			sgin.GetCtx(c).Success("OK")
		})
	})
	ctx.PrintAddr()
	r.Run(ctx.Addr)
}
