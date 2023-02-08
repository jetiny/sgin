package uses

import (
	"jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const appKey = "appKey"

func withAppModel(handler common.AppModelHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		appCode := c.Request.Header.Get(gEnvHeaderAppCode.String())
		if appCode != "" {
			c.Set(appKey, handler(c, appCode))
		} else {
			c.Set(appKey, nil)
		}
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
