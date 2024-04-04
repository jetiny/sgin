package base

import (
	"github.com/jetiny/sgin/utils"

	"github.com/sirupsen/logrus"
)

var (
	//Server
	EnvHost           *utils.EnvGetter = utils.GetterDefault("HOST", "0.0.0.0")
	EnvPort           *utils.EnvGetter = utils.GetterDefault("PORT", 8888)
	EnvNode           *utils.EnvGetter = utils.GetterDefault("Node", 1)
	EnvServerEndpoint *utils.EnvGetter = utils.GetterDefault("SERVER_ENDPOINT", "")

	// IP
	EnvIPUseCwd   *utils.EnvGetter = utils.GetterDefault("IP_USE_CWD", false)
	EnvIPDatabase *utils.EnvGetter = utils.GetterDefault("IP_DATABASE", "./GeoLite2-City.mmdb")

	// Redis
	gEnvRedisAddr     *utils.EnvGetter = utils.GetterDefault("REDIS_ADDR", "0.0.0.0:6379")
	gEnvRedisPassword *utils.EnvGetter = utils.GetterDefault("REDIS_PASSWD", "")
	gEnvRedisDb       *utils.EnvGetter = utils.GetterDefault("REDIS_DB", 0)
	// DB
	gEnvDbEngine  *utils.EnvGetter = utils.GetterDefault("DB_ENGINE", "mysql")
	gEnvDbAddress *utils.EnvGetter = utils.GetterDefault("DB_ADDRESS", "")
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
