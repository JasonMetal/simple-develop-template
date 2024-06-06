package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"time"

	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/number"
)

type RedisInstance struct {
	Client *redis.Client
}

var redisPoolList = make(map[string][]RedisInstance)

// A ping is set to the server with this period to test for the health of
// the connection and server.
const healthCheckPeriod = time.Minute

func GetRedisInstance(dbName string) (*RedisInstance, error) {

	if list, ok := redisPoolList[dbName]; ok {
		l := len(redisPoolList[dbName])
		if l > 0 {
			l = l - 1
			i := number.GetRandNum(l)

			return &list[i], nil
		}
	}

	return nil, errors.New(dbName + " redis list is nil")
}

func SetRedisInstance(dbName string, instance []RedisInstance) {
	redisPoolList[dbName] = instance
}

// Set 用法：Set("key", val, 60)，其中 expire 的单位为秒
func (r *RedisInstance) Set(ctx context.Context, key string, val interface{}, expire int) (string, error) {

	return r.Client.Set(ctx, key, val, time.Second*time.Duration(expire)).Result()
}

func (r *RedisInstance) SetNx(ctx context.Context, key string, val interface{}, expire int) (bool, error) {

	return r.Client.SetNX(ctx, key, val, time.Second*time.Duration(expire)).Result()
}

func (r *RedisInstance) Get(ctx context.Context, key string) (string, error) {

	return r.Client.Get(ctx, key).Result()
}

// SetNxEx exist set value + expires otherwise not do
// cmd: set key value ex 3600 nx
func (r *RedisInstance) SetNxEx(ctx context.Context, key string, val interface{}, expire int) (string, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		value = string(b)
	}
	return r.Client.SetEX(ctx, key, value, time.Second*time.Duration(expire)).Result()
}

func (r *RedisInstance) Del(ctx context.Context, key string) error {
	_, err := r.Client.Del(ctx, key).Result()
	return err
}

func (r *RedisInstance) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

func (r *RedisInstance) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.IncrBy(ctx, key, value).Result()
}

func (r *RedisInstance) Decr(ctx context.Context, key string) (int64, error) {
	return r.Client.Decr(ctx, key).Result()
}

func (r *RedisInstance) DecrBy(ctx context.Context, key string, decrAmount int64) (int64, error) {
	return r.Client.DecrBy(ctx, key, decrAmount).Result()
}

// Lpop 用法：Lpop("key")
func (r *RedisInstance) Lpop(ctx context.Context, key string) (string, error) {

	return r.Client.LPop(ctx, key).Result()
}

func (r *RedisInstance) Llen(ctx context.Context, key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}

// Lpush 用法：Lpush("key", val)
func (r *RedisInstance) Lpush(ctx context.Context, key string, val interface{}) (int64, error) {

	return r.Client.LPush(ctx, key, val).Result()
}

// Rpop 用法：Rpop("key")
func (r *RedisInstance) Rpop(ctx context.Context, key string) (string, error) {
	return r.Client.RPop(ctx, key).Result()
}

func (r *RedisInstance) Lrem(ctx context.Context, key string, val interface{}, count int64) (int64, error) {
	return r.Client.LRem(ctx, key, count, val).Result()
}

// 有序集合操作

func (r *RedisInstance) Zadd(ctx context.Context, key string, score float64, member interface{}) (int64, error) {
	mem := &redis.Z{
		Score:  score,
		Member: member,
	}
	return r.Client.ZAdd(ctx, key, mem).Result()
}

func (r *RedisInstance) Zrangebyscore(ctx context.Context, key, min, max string) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min: min,
		Max: max,
	}
	return r.Client.ZRangeByScore(ctx, key, opt).Result()
}

func (r *RedisInstance) Zrevrange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.ZRevRange(ctx, key, start, stop).Result()
}

func (r *RedisInstance) Hgetall(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *RedisInstance) Hget(ctx context.Context, key, field string) (string, error) {

	return r.Client.HGet(ctx, key, field).Result()
}
func (r *RedisInstance) Hlen(ctx context.Context, key string) (int64, error) {

	return r.Client.HLen(ctx, key).Result()
}

func (r *RedisInstance) Hkeys(ctx context.Context, key string) ([]string, error) {
	return r.Client.HKeys(ctx, key).Result()
}

func (r *RedisInstance) Hset(ctx context.Context, key, field string, val interface{}) (int64, error) {

	return r.Client.HSet(ctx, key, field, val).Result()
}

func (r *RedisInstance) Hdel(ctx context.Context, key, field string) (int64, error) {

	return r.Client.HDel(ctx, key, field).Result()
}

func (r *RedisInstance) Hincrby(ctx context.Context, key, field string, incrAmount int64) (int64, error) {

	return r.Client.HIncrBy(ctx, key, field, incrAmount).Result()
}

// Expire 设置键过期时间，expire的单位为秒
func (r *RedisInstance) Expire(ctx context.Context, key string, expire int) error {
	_, err := r.Client.Expire(ctx, key, time.Second*time.Duration(expire)).Result()
	return err
}

func (r *RedisInstance) GetString(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisInstance) SetString(ctx context.Context, key string, value string) (string, error) {

	return r.Client.Set(ctx, key, value, 0).Result()
}

func (r *RedisInstance) Ttl(ctx context.Context, key string) (int64, error) {
	duration, err := r.Client.TTL(ctx, key).Result()

	return int64(duration), err
}

func (r *RedisInstance) SAdd(ctx context.Context, key string, value string) (int64, error) {
	return r.Client.SAdd(ctx, key, value).Result()
}
func (r *RedisInstance) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}
func (r *RedisInstance) SRANDMEMBER(ctx context.Context, key string) (string, error) {

	return r.Client.SRandMember(ctx, key).Result()
}

func (r *RedisInstance) Pub(ctx context.Context, channel string, msg string) (int64, error) {
	return r.Client.Publish(ctx, channel, msg).Result()
}

func (r *RedisInstance) Sub(ctx context.Context, consumeFunc func(data *redis.Message) error, channel string) error {
	//return r.Do(ctx, "SUBSCRIBE", channel)

	psc := r.Client.Subscribe(ctx, channel)
	if err := psc.Ping(ctx, ""); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {

		for {
			rec, err := psc.ReceiveMessage(ctx)
			if err != nil {
				return
			}
			if err := consumeFunc(rec); err != nil {
				done <- err
				return
			}
		}
	}()
	// test for the health
	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if err := psc.Unsubscribe(ctx); err != nil {
				return fmt.Errorf("redis pubsub unsubscribe err: %v", err)
			}
			return nil
		case err := <-done:
			return err
		case <-ticker.C:
			if err := psc.Ping(ctx); err != nil {

				return err
			}
		}
	}

}
