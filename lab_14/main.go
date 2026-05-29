package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func IsEven(n int) bool {
	return n%2 == 0
}

func main() {
	contacts := []Contact{
		{ID: 1, Name: "Tony Stark", Phone: "099-111-2233"},
		{ID: 2, Name: "Steve Rogers", Phone: "095-222-4455"},
	}

	http.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(contacts)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API is running perfectly!")
	})

	fmt.Println("Сервер запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Помилка запуску сервера: %v\n", err)
	}
}
