package authentication

import (
	"time"

	"github.com/go-redis/redis"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type redisCache struct {
	Host     string
	ID       int
	Password string
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.Host,
		Password: cache.Password,
		DB:       cache.ID,
	})
}

func (cache *redisCache) Get(key string) (string, error) {
	client := cache.getClient()
	result, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (cache *redisCache) Set(key string, value interface{}, exp time.Duration) error {
	client := cache.getClient()
	err := client.Set(key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewCache(host, password string, id int) *redisCache {
	return &redisCache{
		Host:     host,
		ID:       id,
		Password: password,
	}
}
