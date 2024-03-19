package uses

import (
	"github.com/go-xorm/xorm"
	"github.com/jetiny/sgin/common"
)

// type btx struct {
// }

// func (b *btx) Orm() *xorm.Engine {
// 	return common.GetBootContext().Xorm
// }

// func (b *btx) Redis() *rediCache {
// 	return &rediCache{
// 		r: common.GetBootContext().Redis,
// 	}
// }

func Orm() *xorm.Engine {
	return common.GetBootContext().Xorm
}

func Redis() *rediCache {
	return &rediCache{
		r: common.GetBootContext().Redis,
	}
}
