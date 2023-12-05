package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketURL      = getEnv("WS_URL", "wss://api.tiingo.com/fx")
	priceFeedAPIToken = getEnv("WS_API_KEY", "15fdaffbca93fb6c1084fb284f974be97ef23dcf")
	timestampLayout   = getEnv("WS_TIME_LAYOUT", "2006-01-02T15:04:05.000000-07:00")
)

func connectWebSocket() *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		logger.Fatalf("websocket dial error: %v\n", err)
	}

	subscribeRequest, err := json.Marshal(map[string]any{
		"eventName":     "subscribe",
		"authorization": priceFeedAPIToken,
		"eventData": map[string]any{
			"thresholdLevel": 5,
			"tickers":        []string{"gbpusd"},
		},
	})
	if err != nil {
		logger.Fatalf("json marshal error: %v\n", err)
	}

	if err = c.WriteMessage(websocket.TextMessage, subscribeRequest); err != nil {
		logger.Fatalf("write error: %v\n", err)
	}
	return c
}

func pollWebSocket(c *websocket.Conn) {
	defer c.Close()

	latestUpdate, err := retrieveDataRedis()
	if err != nil {
		logger.Printf("redis db read error: %v\n", err)
		latestUpdate = Data{Timestamp: time.Now().UTC()}
	}

	for {
		time.Sleep(time.Second)
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Printf("websocket read error: %v\n", err)
			break
		}

		var wsResponse WebSocketResponse
		if err := json.Unmarshal(message, &wsResponse); err != nil {
			logger.Printf("json unmarshal error: %v\n", err)
			continue
		}

		if !(wsResponse.MessageType == "A" && wsResponse.Service == "fx") {
			continue
		}

		var data []any
		if err := json.Unmarshal(wsResponse.Data, &data); err != nil || len(data) < 6 {
			logger.Printf("invalid fx price data: %s\n", wsResponse.Data)
			continue
		}

		t, ok := data[2].(string)
		if !ok || t == "" {
			logger.Printf("invalid timestamp update: %s\n", t)
			continue
		}

		timestamp, err := time.Parse(timestampLayout, t)
		if err != nil {
			logger.Printf("error parsing timestamp: %s\n", timestamp)
			continue
		}

		timestamp = timestamp.UTC()
		if timestamp.Before(latestUpdate.Timestamp) || timestamp.Equal(latestUpdate.Timestamp) || err != nil {
			logger.Printf("expired timestamp: received %s, cached is %s\n", timestamp, latestUpdate.Timestamp)
			continue
		}

		priceFloat, ok := data[5].(float64)
		if !ok || priceFloat == 0 {
			logger.Printf("invalid price update: %f\n", priceFloat)
			continue
		}

		price := fmt.Sprintf("%f", priceFloat)
		if price == latestUpdate.Value {
			continue
		}

		simplifiedData := Data{
			Timestamp: timestamp,
			Value:     price,
		}

		if err := storeDataRedis(simplifiedData); err != nil {
			logger.Printf("redis write error: %v\n", err)
			continue
		}

		latestUpdate = simplifiedData
		logger.Printf("\nSuccessfully fetched & cached new price update from WebSocket feed:\nGBP/USD @ %s\n\n", price)
	}
}

type WebSocketResponse struct {
	MessageType string          `json:"messageType"`
	Service     string          `json:"service,omitempty"`
	Data        json.RawMessage `json:"data"`
}
