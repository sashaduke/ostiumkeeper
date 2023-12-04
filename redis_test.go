package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestRedis(t *testing.T) {
	data, err := retrieveDataRedis()
	fmt.Println(data, err)
}
