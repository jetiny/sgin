package uses

import (
	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const tokenKey = "tokenKey"

func getUserToken(c *gin.Context) *common.UserToken {
	value := c.MustGet(tokenKey)
	if value != nil {
		return value.(*common.UserToken)
	}
	return nil
}

func hasUserToken(c *gin.Context) *common.UserToken {
	value, exists := c.Get(tokenKey)
	if exists && value != nil {
		return value.(*common.UserToken)
	}
	return nil
}
