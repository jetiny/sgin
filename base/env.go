package base

import (
	"jetiny/sgin/utils"

	"github.com/sirupsen/logrus"
)

const (
	ActionNodeCountOne = 1 // 单节点
	ActionNodeId0      = 0 // 单节点
)

var (
	//Server
	EnvHost   *utils.EnvGetter = utils.GetterDefault("HOST", "0.0.0.0")
	EnvPort   *utils.EnvGetter = utils.GetterDefault("PORT", "8888")
	EnvNode   *utils.EnvGetter = utils.GetterDefault("Node", 1)
	ServerUrl *utils.EnvGetter = utils.GetterDefault("SERVER_URL", "http://localhost:8888")

	DebugToken         *utils.EnvGetter = utils.GetterDefault("DEBUG_TOKEN", "DebugToken")
	TestFlags          *utils.EnvGetter = utils.GetterDefault("TEST_FLAGS", uint64(1))
	EnvAuthPrefixToken *utils.EnvGetter = utils.GetterDefault("PREFIX_TOKEN", "token:")
	// Redis
	EnvRedisAddr     *utils.EnvGetter = utils.GetterDefault("REDIS_ADDR", "0.0.0.0:6379")
	EnvRedisPassword *utils.EnvGetter = utils.GetterDefault("REDIS_PASSWD", "")
	EnvRedisDb       *utils.EnvGetter = utils.GetterDefault("REDIS_DB", 0)
	EnvRedisPoolSize *utils.EnvGetter = utils.GetterDefault("REDIS_POOL_SIZE", 10)
	EnvRedisTimeout  *utils.EnvGetter = utils.GetterDefault("REDIS_TIMEOUT", 5)
	// Mysql
	EnvMysqlAddr *utils.EnvGetter = utils.GetterDefault("MYSQL_ADDR", "root@tcp(127.0.0.1:3306)/jetserver?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local")
	// Session
	EnvSessionSecret       *utils.EnvGetter = utils.GetterDefault("SESSION_SECRET", "ss")
	EnvSessionStorePrefix  *utils.EnvGetter = utils.GetterDefault("SESSION_PREFIX", "ss:")
	EnvSessionStoreMaxSize *utils.EnvGetter = utils.GetterDefault("SESSION_MAXSIZE", 4096)
	EnvSessionKey          *utils.EnvGetter = utils.GetterDefault("SESSION_KEY", "S")
	EnvSessionExpiredHour  *utils.EnvGetter = utils.GetterDefault("SESSION_EXPIRED", 2) // 小时
	// Oss
	OssServer         *utils.EnvGetter = utils.GetterDefault("OSS_SERVER", "https://oss-cn-hangzhou.aliyuncs.com/")
	OssServerInternal *utils.EnvGetter = utils.GetterDefault("OSS_SERVER_INTERNAL", "https://oss-cn-hangzhou.aliyuncs.com/")
	OssRegion         *utils.EnvGetter = utils.GetterDefault("OSS_REGION", "cn-hangzhou")
	OssBucketName     *utils.EnvGetter = utils.Getter("OSS_BUCKET")
	OssAccount        *utils.EnvGetter = utils.Getter("OSS_ACCOUNT")
	OssSecret         *utils.EnvGetter = utils.Getter("OSS_SECRET")
	OssBaseDir        *utils.EnvGetter = utils.Getter("OSS_BASE_DIR")
	StorageType       *utils.EnvGetter = utils.GetterDefault("STORAGE_TYPE", "oss") // oss
	LocalStoragePath  *utils.EnvGetter = utils.GetterDefault("LOCAL_STORAGE_PATH", "storage")
	// Log
	EnvLogDir      *utils.EnvGetter = utils.GetterDefault("LOG_DIR", "logs")
	EnvLogFileName *utils.EnvGetter = utils.GetterDefault("LOG_FILENAME", "app.log")
	EnvLogLevel    *utils.EnvGetter = utils.GetterDefault("LOG_LEVEL", int(logrus.TraceLevel)) // Level
	EnvLogExpired  *utils.EnvGetter = utils.GetterDefault("LOG_EXPIRED", 7)                    // 日志 文件过期时间
	EnvLogCutDays  *utils.EnvGetter = utils.GetterDefault("LOG_CUT_DAYS", 1)                   // 日志切割时间

	EnableMock          *utils.EnvGetter = utils.GetterDefault("ENABLE_MOCK", false)                    // 日志切割时间
	ActionNodeCount     *utils.EnvGetter = utils.GetterDefault("ACTION_NODE_COUNT", ActionNodeCountOne) // action 节点个数
	ActionNodeId        *utils.EnvGetter = utils.GetterDefault("ACTION_NODE_Id", ActionNodeId0)         // action 节点ID
	LibraryPersonalName *utils.EnvGetter = utils.GetterDefault("LIBRARY_PERSONAL_NAME", "Library")      // 文件夹默认前缀
	LibraryFolderPrefix *utils.EnvGetter = utils.GetterDefault("LIBRARY_FOLDER_PREFIX", "$root")        // 文件夹默认前缀
	LibraryTrashName    *utils.EnvGetter = utils.GetterDefault("LIBRARY_TRASH_NAME", "Trash")           // 文件夹默认前缀
	LibraryTrashPrefix  *utils.EnvGetter = utils.GetterDefault("LIBRARY_TRASH_PREFIX", "$trash")        // 文件夹默认前缀
	LibraryRootName     *utils.EnvGetter = utils.GetterDefault("LIBRARY_ROOT_NAME", "Root")             // 文件夹默认前缀
	LibraryRootPrefix   *utils.EnvGetter = utils.GetterDefault("LIBRARY_ROOT_PREFIX", "$root")          // 文件夹默认前缀
	EnableFileCalc      *utils.EnvGetter = utils.GetterDefault("ENABLE_FILE_CALC", true)                // 是否开启文件统计
)
