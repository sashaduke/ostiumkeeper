package main

import (
	"context"
	"encoding/json"
	"time"
)

// Data is a type used to contain the timestamped pricefeed data.
type Data struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

// storeDataRedis caches data in Redis.
func storeDataRedis(data Data) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Fatalf("json marshal error: %v\n", err)
		return err
	}

	if err := rdb.Set(context.Background(), "fxPriceData", jsonData, 0).Err(); err != nil {
		logger.Fatalf("redis set error: %v\n", err)
		return err
	}
	return nil
}

// retrieveDataFromRedis fetches data from Redis.
func retrieveDataRedis() (Data, error) {
	val, err := rdb.Get(context.Background(), "fxPriceData").Result()
	if err != nil {
		return Data{}, err
	}

	var data Data
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Data{}, err
	}
	return data, nil
}
