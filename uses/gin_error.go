package uses

import (
	"jetiny/sgin/utils"

	"github.com/gin-gonic/gin"
)

func withGinError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				var Err *utils.Error
				if e, ok := err.(*utils.Error); ok {
					Err = e
				} else {
					Err = utils.InternalServerError.WithMessage(e.Error())
				}
				c.JSON(Err.StatusCode, Err)
				return
			}
		}
	}
}
