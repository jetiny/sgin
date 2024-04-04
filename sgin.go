package sgin

import (
	"github.com/gin-gonic/gin"
	"github.com/jetiny/sgin/base"
	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/uses"
)

func MustBoot(features base.BootFeature) *common.BootContext {
	return base.MustBoot(features)
}

func MustBootAll() *common.BootContext {
	return base.MustBoot(base.BootWithDefault)
}

func GetCtx(c *gin.Context) *uses.Ctx {
	return uses.GetCtx(c)
}
