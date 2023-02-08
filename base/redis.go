package base

import (
	"context"

	"github.com/jetiny/sgin/common"

	"github.com/go-redis/redis/v8"
)

type redisHook struct{}

func (hook *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	common.Logger.Trace("BeforeProcess", cmd)
	return ctx, nil
}

func (hook *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	common.Logger.Trace("AfterProcess", cmd)
	return nil
}

func (hook *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	common.Logger.Trace("BeforeProcessPipeline", cmds)
	return ctx, nil
}
func (hook *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	common.Logger.Trace("AfterProcessPipeline", cmds)
	return nil
}

func initRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     gEnvRedisAddr.String(),
		Password: gEnvRedisPassword.String(), // no password set
		DB:       gEnvRedisDb.Int(),          // use default DB
		PoolSize: gEnvRedisDb.Int(),          // 连接池大小
	})
	rdb.AddHook(&redisHook{})
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(EnvRedisTimeout.Int())*time.Second)
	// defer cancel()
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
