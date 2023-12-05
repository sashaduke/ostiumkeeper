package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.FlushDB(context.Background()).Err(); err != nil {
		panic(err)
	}
	storeDataRedis(Data{Timestamp: time.Now().UTC()})
}

func main() {
	go pollWebSocket(connectWebSocket())
	go keeper()

	http.HandleFunc("/data", handleData)
	http.HandleFunc("/contracts", handleContracts)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
