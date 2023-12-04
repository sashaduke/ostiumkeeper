package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const PriceFeedAPIToken = "15fdaffbca93fb6c1084fb284f974be97ef23dcf"

// connectWebSocket connects to a WebSocket and handles incoming messages.
func connectWebSocket(url string) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("dial error: %v\n", err)
	}
	defer c.Close()

	subscribe := map[string]any{
		"eventName":     "subscribe",
		"authorization": PriceFeedAPIToken,
		"eventData": map[string]any{
			"thresholdLevel": 5,
			"tickers":        []string{"gbpusd"},
		},
	}

	subscribeJSON, err := json.Marshal(subscribe)
	if err != nil {
		log.Fatalf("json marshal error: %v\n", err)
	}

	err = c.WriteMessage(websocket.TextMessage, subscribeJSON)
	if err != nil {
		log.Fatalf("write error: %v\n", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read error: %v\n", err)
			break
		}

		var data Data
		if err := json.Unmarshal(message, &data); err != nil {
			log.Printf("json unmarshal error: %v\n", err)
			continue
		}
		storeDataRedis(data)

		time.Sleep(100 * time.Millisecond)
	}
}
