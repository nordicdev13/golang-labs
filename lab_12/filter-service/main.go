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
		log.Fatalf("Filter не зміг підключитися до NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Filter Service чекає на числа...")

	_, err = nc.Subscribe("pipeline.numbers", func(m *nats.Msg) {
		var msg Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("Помилка розбору JSON: %v", err)
			return
		}

		if msg.Value%2 == 0 {
			log.Printf("[Filtered] Число %d парне, пересилаємо далі", msg.Value)
			nc.Publish("pipeline.even", m.Data)
		}
	})
	if err != nil {
		log.Fatalf("Помилка підписки: %v", err)
	}

	select {}
}
