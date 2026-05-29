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
		log.Fatalf("Sum не зміг підключитися до NATS: %v", err)
	}
	defer nc.Close()

	var totalSum int
	log.Println("Sum Service готовий рахувати суму...")

	_, err = nc.Subscribe("pipeline.squared", func(m *nats.Msg) {
		var msg Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("Помилка розбору JSON: %v", err)
			return
		}

		totalSum += msg.Value
		log.Printf("[Sum] Додано: %d. Поточна сума всіх чисел: %d", msg.Value, totalSum)
	})
	if err != nil {
		log.Fatalf("Помилка підписки: %v", err)
	}

	select {}
}
