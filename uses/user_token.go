package uses

import (
	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const tokenKey = "tokenKey"

func withUserToken(handler common.TokenHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(gEnvHeaderAccessToken.String())
		if token != "" {
			c.Set(tokenKey, handler(c, token))
		} else {
			c.Set(tokenKey, nil)
		}
		c.Next()
	}
}

func getUserToken(c *gin.Context) *common.UserToken {
	value := c.MustGet(tokenKey)
	if value != nil {
		return value.(*common.UserToken)
	}
	return nil
}
