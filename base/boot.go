package base

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jetiny/sgin/common"
	"github.com/jetiny/sgin/utils"
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
	// 使用Ip
	BootWithIp BootFeature = 1 << 5
	// 默认全部开启
	BootWithDefault BootFeature = BootWithEnv |
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
		Routes:        make([]*common.Route, 0),
		Tasks:         make(map[string][]gin.HandlerFunc),
		TokenHandlers: make(map[string]common.TokenHandler),
		Addr:          EnvHost.String() + ":" + strconv.Itoa(EnvPort.Int()),
	}
	if hasFeature(features, BootWithEnv) {
		err := godotenv.Load()
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, bootError("env", err)
			}
		}
	}
	if hasFeature(features, BootWithLogger) {
		err := initLogger()
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
		engine, err := initDb()
		if err != nil {
			return nil, bootError("mysql", err)
		}
		if hasFeature(features, BootWithLogger) {
			engine.SetLogger(gxormLogger)
		}
		res.Xorm = engine
	}
	if hasFeature(features, BootWithIp) {
		ip, err := initIp()
		if err != nil {
			return nil, bootError("ip", err)
		}
		res.IpEngine = ip
	}
	common.SetBootContext(res)
	return res, nil
}

func MustBoot(features BootFeature) *common.BootContext {
	resource, err := boot(features)
	if err != nil {
		log.Fatal(err)
	}
	return resource
}
