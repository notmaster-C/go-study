package cache

import (
	"context"
	"fmt"
	"go-study/config"
	logging "go-study/log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Nil = redis.Nil
	ctx = context.Background()
	log = logging.Logger("cache")
	// Redis redis客户端
	Redis *redis.Client
)

const prefix = "go_study"

// 初始化redis连接
func InitRedis(cfg *config.Redis) {
	host := fmt.Sprintf("%v:%v", cfg.Addr, cfg.Port)

	Redis = redis.NewClient(&redis.Options{
		Addr:       host,
		Password:   cfg.Passwd,
		DB:         cfg.Db,
		PoolSize:   cfg.PoolSize,
		MaxRetries: cfg.MaxRetries,
	})

	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Panicf("failed to create redis client error:%v", err)
	}
	fmt.Println("successfully connect redis")
	// log.Infof("successfully connect redis [%v:%v/%d] maxRetries:%d poolSize:%d", host, cfg.Passwd, cfg.Db, cfg.MaxRetries, cfg.PoolSize)
}

func RCli() *redis.Client {
	return Redis
}

func Incr(key string) error {
	key = fmt.Sprintf("%v.%v", prefix, key)
	if err := Redis.Incr(ctx, key).Err(); err != nil {
		return Redis.Expire(ctx, key, time.Hour).Err()
	}
	return nil
}

// 增加前缀
func Set(key string, value interface{}, expiration time.Duration) (string, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	return Redis.Set(ctx, key, value, expiration).Result()
}

func Get(key string) (string, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	return Redis.Get(ctx, key).Result()
}
func Del(key string) (int64, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	return Redis.Del(ctx, key).Result()
}
func HSet(key, field string, value interface{}, expiration time.Duration) error {
	key = fmt.Sprintf("%v.%v", prefix, key)
	err := Redis.HSet(ctx, key, field, value).Err()
	if err == redis.Nil {
		err = nil
	}
	Redis.Expire(ctx, key, expiration).Result()
	if err != nil {
		log.Warningf("faild to redis hset key[%v] field[%v] value[%v] error: %v", key, field, value, err)
	}
	return err
}

func HGet(key, field string) (string, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	res, err := Redis.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		log.Warningf("faild to redis hget key[%v] field[%v] value[%v] error: %v", key, field, err)
	}
	return res, err
}

func HDel(key, field string) error {
	key = fmt.Sprintf("%v.%v", prefix, key)
	_, err := Redis.HDel(ctx, key, field).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		log.Warningf("faild to redis hdel key[%v] field[%v] value[%v] error: %v", key, field, err)
	}
	return err
}

func HKeys(key string) ([]string, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	res, err := Redis.HKeys(ctx, key).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		log.Warningf("faild to redis hget key[%v] field[%v] value[%v] error: %v", key, err)
	}
	return res, err
}

func HLen(key string) (int64, error) {
	key = fmt.Sprintf("%v.%v", prefix, key)
	res, err := Redis.HLen(ctx, key).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		log.Warningf("faild to redis hget key[%v] field[%v] value[%v] error: %v", key, err)
	}
	return res, err
}

func IsExist(key string) bool {
	key = fmt.Sprintf("%v.%v", prefix, key)
	res := Redis.Exists(ctx, key).Val()
	if res == 0 {
		return false
	}
	return true
}
