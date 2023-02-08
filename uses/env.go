package uses

import "github.com/jetiny/sgin/utils"

var (
	// Redis
	gEnvRedisKeyPrefix   *utils.EnvGetter = utils.GetterDefault("REDIS_KEY_PREFIX", "app:")
	gEnvRedisTokenPrefix *utils.EnvGetter = utils.GetterDefault("REDIS_TOKEN_PREFIX", "token:")
	// Auth
	gEnvHeaderAppCode     *utils.EnvGetter = utils.GetterDefault("HEADER_APP_CODE", "App-Code")
	gEnvHeaderAccessToken *utils.EnvGetter = utils.GetterDefault("HEADER_ACCESS_TOKEN", "Access-Token")
	// Cors
	gEnvCorsAllowOrigins  = utils.GetterDefault("CORS_ALLOW_ORIGINS", "*")
	gEnvCorsAllowMethods  = utils.GetterDefault("CORS_ALLOW_METHODS", "*")
	gEnvCorsAllowHeaders  = utils.GetterDefault("CORS_ALLOW_HEADERS", "*")
	gEnvCorsExposeHeaders = utils.GetterDefault("CORS_EXPOSE_HEADERS", "Content-Length")
)
