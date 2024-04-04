package uses

import (
	"github.com/go-xorm/xorm"
	"github.com/jetiny/sgin/common"
	"github.com/oschwald/geoip2-golang"
)

func Orm() *xorm.Engine {
	return common.GetBootContext().Xorm
}

func Redis() *rediCache {
	return &rediCache{
		r: common.GetBootContext().Redis,
	}
}

func IpEngine() *geoip2.Reader {
	return common.GetBootContext().IpEngine
}
