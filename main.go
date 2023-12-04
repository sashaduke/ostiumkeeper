package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func main() {
	go pollWebSocket(connectWebSocket("wss://api.tiingo.com/fx"))
	go keeper()

	http.HandleFunc("/data", handleData)
	http.HandleFunc("/contracts", handleContracts)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
