package uses

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const redisKey = "redisKey"

func withRedis(engine *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(redisKey, &rediCache{r: engine})
		c.Next()
	}
}

func getRedis(c *gin.Context) *rediCache {
	value := c.MustGet(redisKey)
	return value.(*rediCache)
}

type rediCache struct {
	r *redis.Client
}

func (rc rediCache) Client() *redis.Client {
	return rc.r
}

func (rc *rediCache) Set(key string, value any, dur time.Duration) error {
	key = gEnvRedisKeyPrefix.String() + key
	return rc.r.Set(context.Background(), key, value, dur).Err()
}

func (rc *rediCache) SetJson(key string, value any, dur time.Duration) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key = gEnvRedisKeyPrefix.String() + key
	value = string(buf)
	return rc.r.Set(context.Background(), key, value, dur).Err()
}

func (rc *rediCache) GetJson(key string, value any) (bool, error) {
	var bytes []byte
	key = gEnvRedisKeyPrefix.String() + key
	str, err := rc.r.Get(context.Background(), key).Bytes()
	if isRedisNil(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	bytes = str
	return true, json.Unmarshal(bytes, value)
}

func (rc *rediCache) GetInt(key string) (int, bool) {
	key = gEnvRedisKeyPrefix.String() + key
	value, err := rc.r.Get(context.Background(), key).Int()
	if isRedisNil(err) {
		return value, false
	}
	return value, true
}

func (rc *rediCache) GetInt64(key string) (int64, bool) {
	key = gEnvRedisKeyPrefix.String() + key
	value, err := rc.r.Get(context.Background(), key).Int64()
	if isRedisNil(err) {
		return value, false
	}
	return value, true
}

func (rc *rediCache) GetUint64(key string) (uint64, bool) {
	key = gEnvRedisKeyPrefix.String() + key
	value, err := rc.r.Get(context.Background(), key).Uint64()
	if isRedisNil(err) {
		return value, false
	}
	return value, true
}

func (rc *rediCache) GetBool(key string) (bool, bool) {
	key = gEnvRedisKeyPrefix.String() + key
	value, err := rc.r.Get(context.Background(), key).Bool()
	if isRedisNil(err) {
		return value, false
	}
	return value, true
}

func (rc *rediCache) GetString(key string) string {
	key = gEnvRedisKeyPrefix.String() + key
	return rc.r.Get(context.Background(), key).Val()
}

func (rc *rediCache) SetToken(token *common.UserToken) error {
	return rc.SetJson(gEnvRedisTokenPrefix.String()+token.AccessToken, token, time.Until(token.AccessExpiredAt))
}

func (rc *rediCache) GetToken(accessToken string) (*common.UserToken, error) {
	var r common.UserToken
	exists, err := rc.GetJson(gEnvRedisTokenPrefix.String()+accessToken, &r)
	if err == nil && exists {
		return &r, nil
	}
	return nil, err
}

func (rc *rediCache) Del(key string) error {
	err := rc.r.Del(context.Background(), gEnvRedisTokenPrefix.String()+key).Err()
	if isRedisNil(err) {
		return nil
	}
	return err
}

func (rc *rediCache) DelToken(accessToken string) error {
	err := rc.r.Del(context.Background(), gEnvRedisTokenPrefix.String()+accessToken).Err()
	if isRedisNil(err) {
		return nil
	}
	return err
}

func isRedisNil(err error) bool {
	return err == redis.Nil
}
