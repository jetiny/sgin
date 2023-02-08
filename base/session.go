package base

import (
	redisStore "github.com/gin-contrib/sessions/redis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func initSession() (gin.HandlerFunc, error) {
	rstore, err := redisStore.NewStore(10, "tcp", gEnvRedisAddr.String(), "", []byte(gEnvSessionSecret.String()))
	if err != nil {
		return nil, err
	}
	err, realStore := redisStore.GetRedisStore(rstore)
	if err != nil {
		return nil, err
	}
	realStore.DefaultMaxAge = gEnvSessionExpiredHour.Int() * 60 * 60
	realStore.SetKeyPrefix(gEnvSessionStorePrefix.String())
	realStore.SetMaxLength(gEnvSessionStoreMaxSize.Int())
	fun := sessions.Sessions(gEnvSessionKey.String(), rstore)
	return fun, nil
}
