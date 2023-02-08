package uses

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

const ormKey = "ormKey"

func withOrm(engine *xorm.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ormKey, engine)
		c.Next()
	}
}

func getOrm(c *gin.Context) *xorm.Engine {
	value := c.MustGet(ormKey)
	return value.(*xorm.Engine)
}
