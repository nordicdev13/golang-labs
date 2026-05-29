package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

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
		log.Fatalf("Generator не зміг підключитися до NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Generator Service запущено. Починаємо генерацію...")

	time.Sleep(2 * time.Second)

	for i := 1; i <= 100; i++ {
		msg := Message{Value: i}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Помилка маршалінгу: %v", err)
			continue
		}

		err = nc.Publish("pipeline.numbers", jsonData)
		if err != nil {
			log.Printf("Помилка публікації: %v", err)
		} else {
			log.Printf("[Generated] Надіслано число: %d", i)
		}

		time.Sleep(50 * time.Millisecond)
	}

	select {}
}
