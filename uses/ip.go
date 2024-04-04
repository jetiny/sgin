package uses

import (
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

const ipKey = "ipKey"

func withIp(engine *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ipKey, engine)
		c.Next()
	}
}

func getIp(c *gin.Context) *geoip2.Reader {
	value := c.MustGet(ipKey)
	return value.(*geoip2.Reader)
}
