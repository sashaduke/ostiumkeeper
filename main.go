package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func main() {
	// Flush & init DB
	initRedis()

	// Start daemons
	go pollWebSocket(connectWebSocket())
	go keeper()

	// Setup routes
	http.HandleFunc("/data", handleData)
	http.HandleFunc("/contracts", handleContracts)

	// Run server and log if error returned
	logger.Fatal(http.ListenAndServe(":8080", nil))
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.FlushDB(context.Background()).Err(); err != nil {
		logger.Printf("couldn't flush redis DB: %v\n", err)
	}
	if err := storeDataRedis(Data{Timestamp: time.Now().UTC()}); err != nil {
		logger.Printf("redis write error: %v\n", err)
	}
}
