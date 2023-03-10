package uses

import (
	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const appKey = "appKey"

func withAppModel(handler common.AppModelHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		appCode := c.Request.Header.Get(gEnvHeaderAppCode.String())
		c.Set(appKey, handler(c, appCode))
		c.Next()
	}
}

func getAppModel(c *gin.Context) *common.AppModel {
	value := c.MustGet(appKey)
	if value != nil {
		return value.(*common.AppModel)
	}
	return nil
}
