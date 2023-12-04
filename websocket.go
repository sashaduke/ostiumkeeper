package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const PriceFeedAPIToken = "15fdaffbca93fb6c1084fb284f974be97ef23dcf"

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

		var wsResponse WebSocketResponse
		if err := json.Unmarshal(message, &wsResponse); err != nil {
			log.Printf("json unmarshal error: %v\n", err)
			continue
		}

		if wsResponse.MessageType == "A" && wsResponse.Service == "fx" {
			var data []any
			if err := json.Unmarshal(wsResponse.Data, &data); err != nil {
				log.Printf("fx data unmarshal error: %v\n", err)
				continue
			}

			if len(data) > 5 {
				timeStamp, ok := data[2].(string)
				if !ok {
					log.Printf("unexpected data type for timestamp\n")
					continue
				}
				log.Printf("Received timestamp: %s\n", timeStamp)

				pricePoint, ok := data[5].(float64)
				if !ok {
					log.Printf("unexpected data type for price point\n")
					continue
				}
				log.Printf("Received price point: %f\n", pricePoint)

				simplifiedData := Data{
					Timestamp: timeStamp,
					Value:     pricePoint,
				}
				storeDataRedis(simplifiedData)
				fmt.Println(simplifiedData)
			}
		}
		time.Sleep(time.Second)
	}
}

type WebSocketResponse struct {
	MessageType string          `json:"messageType"`
	Service     string          `json:"service,omitempty"`
	Data        json.RawMessage `json:"data"`
}
