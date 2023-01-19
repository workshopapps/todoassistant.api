package queues

import (
	"github.com/go-redis/redis/v9"
)

type RedisQueue struct {
	Client *redis.Client
}

func NewRedisQueue(addr string) (*RedisQueue, error) {
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisQueue{Client: rdb}, nil
}
