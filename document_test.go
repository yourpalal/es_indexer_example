package main

import (
	"encoding/json"
	"testing"
)

// Test_DocumentJsonRoundTrip ensures that the Document class can
// survive a json round-trip.
// uses example_doc from elasticsearch_unit_test.go
func Test_DocumentJsonRoundTrip(t *testing.T) {
	jsonified, err := json.Marshal(example_doc)
	if err != nil {
		t.Fatalf("Error marshaling json: %s", err.Error())
	}

	var doc_copy *Document
	err = json.Unmarshal(jsonified, &doc_copy)
	if err != nil {
		t.Fatalf("Error unmarshaling json: %s", err.Error())
	}

	if doc_copy == nil || example_doc != *doc_copy {
		t.Error("JSON roundtrip produced unlike values")
	}
}
