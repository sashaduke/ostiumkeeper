package main

import (
	"context"
	"encoding/json"
	"log"
)

// Data structure to store the data.
type Data struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// storeDataRedis caches data in Redis.
func storeDataRedis(data Data) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("json marshal error: %v\n", err)
	}

	err = rdb.Set(context.Background(), "fxPriceData", jsonData, 0).Err()
	if err != nil {
		log.Fatalf("redis set error: %v\n", err)
	}
}

// retrieveDataFromRedis fetches data from Redis.
func retrieveDataRedis() (Data, error) {
	val, err := rdb.Get(context.Background(), "fxPriceData").Result()
	if err != nil {
		return Data{}, err
	}

	var data Data
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		return Data{}, err
	}
	return data, nil
}
