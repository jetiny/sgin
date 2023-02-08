package base

import (
	redisStore "github.com/gin-contrib/sessions/redis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func initSession() (gin.HandlerFunc, error) {
	rstore, err := redisStore.NewStore(10, "tcp", EnvRedisAddr.String(), "", []byte(EnvSessionSecret.String()))
	if err != nil {
		return nil, err
	}
	err, realStore := redisStore.GetRedisStore(rstore)
	if err != nil {
		return nil, err
	}
	realStore.DefaultMaxAge = EnvSessionExpiredHour.Int() * 60 * 60
	realStore.SetKeyPrefix(EnvSessionStorePrefix.String())
	realStore.SetMaxLength(EnvSessionStoreMaxSize.Int())
	fun := sessions.Sessions(EnvSessionKey.String(), rstore)
	return fun, nil
}
