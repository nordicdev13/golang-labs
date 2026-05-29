package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Value int `json:"value"`
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Square не зміг підключитися до NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Square Service чекає на парні числа...")

	_, err = nc.Subscribe("pipeline.even", func(m *nats.Msg) {
		var msg Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("Помилка розбору JSON: %v", err)
			return
		}

		squaredValue := msg.Value * msg.Value
		log.Printf("[Squared] %d -> %d", msg.Value, squaredValue)

		newMsg := Message{Value: squaredValue}
		jsonData, _ := json.Marshal(newMsg)
		nc.Publish("pipeline.squared", jsonData)
	})
	if err != nil {
		log.Fatalf("Помилка підписки: %v", err)
	}

	select {}
}
