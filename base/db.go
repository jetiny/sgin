package base

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func initDb() (*xorm.Engine, error) {
	engineType := gEnvDbEngine.String()
	engineAddr := gEnvDbAddress.String()
	if engineType == "mysql" {
		if engineAddr == "" {
			engineAddr = gEnvMysqlAddr.String()
		}
	}
	engine, err := xorm.NewEngine(engineType, engineAddr)
	if err != nil {
		return nil, err
	}
	engine.DatabaseTZ = time.Local // 必须
	engine.TZLocation = time.Local // 必须
	engine.ShowSQL(true)
	return engine, nil
}
