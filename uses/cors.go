package uses

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func withCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(gEnvCorsAllowOrigins.String(), ","),
		AllowMethods:     strings.Split(gEnvCorsAllowMethods.String(), ","),
		AllowHeaders:     strings.Split(gEnvCorsAllowHeaders.String(), ","),
		ExposeHeaders:    strings.Split(gEnvCorsExposeHeaders.String(), ","),
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
