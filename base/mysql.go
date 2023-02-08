package base

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func initMysql() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", gEnvMysqlAddr.String())
	if err != nil {
		return nil, err
	}
	engine.DatabaseTZ = time.Local // 必须
	engine.TZLocation = time.Local // 必须
	engine.ShowSQL(true)
	return engine, nil
}
