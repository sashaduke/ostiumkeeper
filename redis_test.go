package main

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestRedisStoreAndRetrieveData(t *testing.T) {
	tests := []struct {
		name string
		data Data
	}{
		{"EmptyData", Data{}},
		{"ValidData", Data{Timestamp: time.Now().UTC(), Value: "0.12618"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storeDataRedis(tt.data)
			require.Nil(t, err)
			retrievedData, err := retrieveDataRedis()
			require.Nil(t, err)
			require.Equal(t, tt.data.Timestamp, retrievedData.Timestamp)
			require.Equal(t, tt.data.Value, retrievedData.Value)
		})
	}
}
