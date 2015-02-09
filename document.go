package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Timestamp struct {
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Document struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Timestamp Timestamp `json:"timestamp"`
}

func main() {
	doc := Document{
		Title: "Trumpet.ca Programming Problem",
		Body:  "You’ll implement the code for a search indexing system for models such as this...",
		Timestamp: Timestamp{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
	}

	b, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling json: %s", err.Error())
	}

	fmt.Println(string(b))
}