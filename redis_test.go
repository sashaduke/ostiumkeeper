package main

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestRedisGetCachedValue(t *testing.T) {
	data, err := retrieveDataRedis()
	fmt.Println(data, err)
}
