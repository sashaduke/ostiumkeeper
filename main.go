package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func init() { // Executes first, initialises db
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.FlushDB(context.Background()).Err(); err != nil {
		panic(err)
	}
	storeDataRedis(Data{Timestamp: time.Now().UTC()})
}

func main() {
	// Start daemons
	go pollWebSocket(connectWebSocket())
	go keeper()

	// Setup routes
	http.HandleFunc("/data", handleData)
	http.HandleFunc("/contracts", handleContracts)

	// Run server and log if error returned
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
