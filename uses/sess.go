package uses

import (
	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

func GetSess[T any](c *gin.Context) func(fn func(sess *xorm.Session) T) T {
	return func(fn func(sess *xorm.Session) T) T {
		return Sess(getOrm(c), fn)
	}
}

func Sess[T any](orm *xorm.Engine, fn func(sess *xorm.Session) T) T {
	sess := orm.NewSession()
	defer sess.Close()
	err := sess.Begin()
	utils.EnsureNoError(err)
	defer func() {
		if re := recover(); re != nil {
			e := sess.Rollback()
			if e != nil {
				common.Logger.Warn("runWithSess.Rollback", e)
			} else {
				common.Logger.Info("runWithSess", re)
			}
			panic(re)
		}
	}()
	res := fn(sess)
	err = sess.Commit()
	if err != nil {
		e := sess.Rollback()
		if e != nil {
			common.Logger.Warn("runWithSess.Rollback", e)
		}
	}
	return res
}
