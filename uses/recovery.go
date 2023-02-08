package uses

import (
	"encoding/json"
	"log"
	"net/http/httputil"

	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

const stackKey = "stackKey"

type recoveryStack struct {
	stacks []any
}

func (s *recoveryStack) Push(value ...any) {
	s.stacks = append(s.stacks, value...)
}

func getErrorStack(c *gin.Context) *recoveryStack {
	value := c.MustGet(stackKey)
	return value.(*recoveryStack)
}

func newRecoveryStack() *recoveryStack {
	return &recoveryStack{
		stacks: make([]any, 0),
	}
}

func withRecovery() gin.HandlerFunc {
	logger := log.New(gin.DefaultErrorWriter, "\n\n\x1b[31m", log.LstdFlags)
	return func(c *gin.Context) {
		c.Set(stackKey, newRecoveryStack())
		defer func() {
			if err := recover(); err != nil {
				stacks, _ := json.Marshal(getErrorStack(c).stacks)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				reset := string([]byte{27, 91, 48, 109})
				if gin.IsDebugging() {
					logger.Printf("[Recovery] panic recovered:\n\n%s%s\n\n%s\n\n%s\n\n%s", httprequest, goErr.Error(), stacks, goErr.Stack(), reset)
				} else {
					common.Logger.Errorf("[Recovery] panic recovered:\n\n%s%s\n\n%s\n\n%s\n\n%s", httprequest, goErr.Error(), stacks, goErr.Stack(), reset)
				}
				{
					e, ok := err.(utils.Error)
					if ok {
						c.AbortWithStatusJSON(e.JsonError())
						return
					}
				}
				e, ok := err.(*utils.Error)
				if ok {
					c.AbortWithStatusJSON(e.JsonError())
					return
				}
				c.AbortWithStatusJSON(utils.InternalServerError.JsonError())
			}
		}()
		c.Next()
	}
}
