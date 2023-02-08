package base

import (
	"errors"
	"jetiny/sgin/common"
	"jetiny/sgin/utils"
	"log"

	"github.com/joho/godotenv"
)

type BootFeature int

const (
	// 不使用
	BootWithNone BootFeature = 0
	// 使用dotenv
	BootWithEnv BootFeature = 1 << 1
	// 使用session
	BootWithSession BootFeature = 1 << 1
	// 使用redis
	BootWithRedis BootFeature = 1 << 2
	// 使用orm, mysql
	BootWithOrm BootFeature = 1 << 3
	// 使用Logger
	BootWithLogger BootFeature = 1 << 4
	// 默认全部开启
	BootWithAll BootFeature = BootWithEnv |
		BootWithSession |
		BootWithRedis |
		BootWithOrm |
		BootWithLogger
)

func hasFeature(features BootFeature, feature BootFeature) bool {
	return (features & feature) == feature
}

func bootError(name string, err error) error {
	return errors.New("boot error " + name + ":" + err.Error())
}

func boot(features BootFeature) (*common.BootContext, error) {
	res := &common.BootContext{
		Routes: make([]*common.Route, 0),
	}
	if hasFeature(features, BootWithEnv) {
		err := godotenv.Load()
		if err != nil {
			return nil, bootError("env", err)
		}
	}
	if hasFeature(features, BootWithLogger) {
		err := InitLogger()
		if err != nil {
			return nil, bootError("logger", err)
		}
		res.Logger = common.Logger
	}
	err := utils.InitSnowflake(int64(EnvNode.Int()))
	if err != nil {
		return nil, bootError("snowflake", err)
	}
	if hasFeature(features, BootWithRedis) {
		rds, err := initRedis()
		if err != nil {
			return nil, bootError("redis", err)
		}
		if hasFeature(features, BootWithSession) {
			sess, err := initSession()
			if err != nil {
				return nil, bootError("session", err)
			}
			res.SessionHandle = sess
		}
		res.Redis = rds
	}
	if hasFeature(features, BootWithOrm) {
		engine, err := initMysql()
		if err != nil {
			return nil, bootError("mysql", err)
		}
		if hasFeature(features, BootWithLogger) {
			engine.SetLogger(gxormLogger)
		}
	}
	common.SetBootContext(res)
	return res, nil
}

func MustBootAll() *common.BootContext {
	return MustBoot(BootWithAll)
}

func MustBoot(features BootFeature) *common.BootContext {
	resource, err := boot(features)
	if err != nil {
		log.Fatal(err)
	}
	return resource
}
