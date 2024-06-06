package bootstrap

import (
	"fmt"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/config"
	redisHelper "github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/redis"
	redis "github.com/go-redis/redis/v8"
	yCfg "github.com/olebedev/config"
	"os"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedis() {
	dbList := getDbNames("redis")
	for _, dbname := range dbList {
		instances, err := initRedisPool(dbname)
		if err == nil {
			redisHelper.SetRedisInstance(dbname, instances)
		}
	}

	//go closePool()
}

func initRedisPool(dbName string) ([]redisHelper.RedisInstance, error) {
	path := fmt.Sprintf("%sconfig/%s/redis.yml", ProjectPath(), DevEnv)

	cfg, err := config.GetConfig(path)
	if err != nil {
		return nil, err
	}

	//maxIdle, _ := cfg.Int("redis." + dbName + ".maxIdle")
	maxActive, _ := cfg.Int("redis." + dbName + ".maxActive")
	idleTimeout, _ := cfg.Int("redis." + dbName + ".idleTimeout")

	servers, err := cfg.List("redis." + dbName + ".servers")
	if err != nil {
		return nil, err
	}
	redisPools := make([]redisHelper.RedisInstance, len(servers))
	for k, v := range servers {
		address, _ := yCfg.Get(v, "address")
		passwd, _ := yCfg.Get(v, "passwd")
		db, _ := yCfg.Get(v, "db")

		if address == nil {
			os.Exit(1)
		}

		// 建立连接
		rds := &RedisClient{}
		// 使用默认的 context

		// 使用 redis 库里的 NewClient 初始化连接
		rds.Client = redis.NewClient(&redis.Options{
			Addr:        address.(string),
			Password:    passwd.(string),
			DB:          db.(int),
			PoolSize:    maxActive,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
		})
		redisPools[k] = redisHelper.RedisInstance{
			Client: rds.Client,
		}

	}

	return redisPools, nil
}
