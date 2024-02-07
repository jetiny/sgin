package base

import (
	"github.com/jetiny/sgin/utils"

	"github.com/sirupsen/logrus"
)

var (
	//Server
	gEnvHost *utils.EnvGetter = utils.GetterDefault("HOST", "0.0.0.0")
	gEnvPort *utils.EnvGetter = utils.GetterDefault("PORT", 8888)
	gEnvNode *utils.EnvGetter = utils.GetterDefault("Node", 1)

	// Redis
	gEnvRedisAddr     *utils.EnvGetter = utils.GetterDefault("REDIS_ADDR", "0.0.0.0:6379")
	gEnvRedisPassword *utils.EnvGetter = utils.GetterDefault("REDIS_PASSWD", "")
	gEnvRedisDb       *utils.EnvGetter = utils.GetterDefault("REDIS_DB", 0)
	// Mysql
	gEnvDbEngine  *utils.EnvGetter = utils.GetterDefault("DV_ENGINE", "")
	gEnvDbAddress *utils.EnvGetter = utils.GetterDefault("DV_ADDRESS", "")
	gEnvMysqlAddr *utils.EnvGetter = utils.GetterDefault("MYSQL_ADDR", "root@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local")
	// Session
	gEnvSessionSecret       *utils.EnvGetter = utils.GetterDefault("SESSION_SECRET", "ss")
	gEnvSessionStorePrefix  *utils.EnvGetter = utils.GetterDefault("SESSION_PREFIX", "ss:")
	gEnvSessionStoreMaxSize *utils.EnvGetter = utils.GetterDefault("SESSION_MAXSIZE", 4096)
	gEnvSessionKey          *utils.EnvGetter = utils.GetterDefault("SESSION_KEY", "S")
	gEnvSessionExpiredHour  *utils.EnvGetter = utils.GetterDefault("SESSION_EXPIRED", 2) // 小时
	// Log
	gEnvLogDir      *utils.EnvGetter = utils.GetterDefault("LOG_DIR", "logs")
	gEnvLogFileName *utils.EnvGetter = utils.GetterDefault("LOG_FILENAME", "app.log")
	gEnvLogLevel    *utils.EnvGetter = utils.GetterDefault("LOG_LEVEL", int(logrus.TraceLevel)) // Level
	gEnvLogExpired  *utils.EnvGetter = utils.GetterDefault("LOG_EXPIRED", 7)                    // 日志 文件过期时间
	gEnvLogCutDays  *utils.EnvGetter = utils.GetterDefault("LOG_CUT_DAYS", 1)                   // 日志切割时间
)
