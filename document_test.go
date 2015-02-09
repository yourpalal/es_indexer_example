package main

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_DocumentJsonRoundTrip(t *testing.T) {
	doc := Document{
		Title: "Trumpet.ca Programming Problem",
		Body:  "You’ll implement the code for a search indexing system for models such as this...",
		Timestamp: Timestamp{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
	}

	jsonified, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Error marshaling json: %s", err.Error())
	}

	var doc_copy *Document
	err = json.Unmarshal(jsonified, &doc_copy)
	if err != nil {
		t.Fatalf("Error unmarshaling json: %s", err.Error())
	}

	if doc_copy == nil || doc != *doc_copy {
		t.Error("JSON roundtrip produced unlike values")
	}
}
