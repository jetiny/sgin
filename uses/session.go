package uses

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	return session
}
