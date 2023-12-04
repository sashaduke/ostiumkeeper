package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	PriceFeedAPIToken = "15fdaffbca93fb6c1084fb284f974be97ef23dcf"
	timestampLayout   = "2006-01-02T15:04:05.000000-07:00"
)

func connectWebSocket(url string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("dial error: %v\n", err)
	}

	subscribeRequest, err := json.Marshal(map[string]any{
		"eventName":     "subscribe",
		"authorization": PriceFeedAPIToken,
		"eventData": map[string]any{
			"thresholdLevel": 5,
			"tickers":        []string{"gbpusd"},
		},
	})
	if err != nil {
		log.Fatalf("json marshal error: %v\n", err)
	}

	if err = c.WriteMessage(websocket.TextMessage, subscribeRequest); err != nil {
		log.Fatalf("write error: %v\n", err)
	}

	return c
}

func pollWebSocket(c *websocket.Conn) {
	defer c.Close()

	latestUpdate, err := retrieveDataRedis()
	if err != nil {
		log.Fatalf("redis db read error: %v\n", err)
	}

	for {
		time.Sleep(time.Second)
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read error: %v\n", err)
			break
		}

		var wsResponse WebSocketResponse
		if err := json.Unmarshal(message, &wsResponse); err != nil {
			log.Printf("json unmarshal error: %v\n", err)
			continue
		}

		if !(wsResponse.MessageType == "A" && wsResponse.Service == "fx") {
			continue
		}

		var data []any
		if err := json.Unmarshal(wsResponse.Data, &data); err != nil || len(data) < 6 {
			log.Printf("invalid fx data: %s\n", wsResponse.Data)
			continue
		}

		t, ok := data[2].(string)
		if !ok || t == "" {
			log.Printf("invalid timestamp: %s\n", t)
			continue
		}

		timestamp, err := time.Parse(timestampLayout, t)
		if timestamp.Before(latestUpdate.Timestamp) || err != nil {
			log.Printf("invalid timestamp: %s\n", timestamp)
			continue
		}

		price, ok := data[5].(float64)
		if !ok || price == 0 {
			log.Printf("invalid price: %f\n", price)
			continue
		}

		simplifiedData := Data{
			Timestamp: timestamp,
			Value:     price,
		}

		storeDataRedis(simplifiedData)
		latestUpdate = simplifiedData
		log.Printf("\nSuccessfully cached new update from WebSocket:\nGBP/USD @ %f\n\n", price)
	}
}

type WebSocketResponse struct {
	MessageType string          `json:"messageType"`
	Service     string          `json:"service,omitempty"`
	Data        json.RawMessage `json:"data"`
}
