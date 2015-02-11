package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"
)

// unjsonable is a map[int]int that cannot be jsonified
var (
	unjsonable map[int]int
	indexer    *ElasticSearchIndexer
	httpMock   MockPoster

	example_doc = Document{
		Title: "Trumpet.ca Programming Problem",
		Body:  "You’ll implement the code for a search indexing system for models such as this...",
		Timestamp: Timestamp{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
	}
)

// setup to run before each test
func setup_es_unit_test(t *testing.T) {

	unjsonable = make(map[int]int)
	unjsonable[1] = 1

	httpMock = MockPoster{}
	httpMock.Result.StatusCode = 200
	httpMock.Result.Status = "200 OK"
	indexer = &ElasticSearchIndexer{"127.0.0.1", "9200", &httpMock}
}

func Test_ESIndexerReturnsJsonErrors(t *testing.T) {
	setup_es_unit_test(t)

	r, err := indexer.Index("testing", "impossible", "0", true, unjsonable)

	if err == nil {
		t.Fatal("Failed to return json error")
	}
	if r.Created {
		t.Fatal("Set Created to true on failure")
	}
}

func Test_ESIndexer_docURL(t *testing.T) {
	setup_es_unit_test(t)

	indexer.host = "127.0.0.1"
	indexer.port = "9200"
	url := indexer.docURL("twitter", "tweet", "best_tweet")
	expected := "http://127.0.0.1:9200/twitter/tweet/best_tweet"
	AssertEqual(t, url == expected, "Failed to build document URL.", url, expected)

	if url != expected {
		t.Fatalf("Failed to build document URL. expected %s, got %s", expected, url)
	}
}

// Test_ESIndexer_create_http_errors tests whether the ElasticSearchIndexer
// will properly propogate http errors
func Test_ESIndexer_create_http_errors(t *testing.T) {
	setup_es_unit_test(t)
	httpMock.Err = errors.New("test error")
	response, err := indexer.Index("trumpet", "doc", "first", true, example_doc)

	AssertEqual(t, err == httpMock.Err, "ElasticSearchIndexer did not propagate error correctly", err, httpMock.Err)
	AssertFalse(t, "indexer reports create on error?", response.Created)
}

func Test_ESIndexer_create(t *testing.T) {
	setup_es_unit_test(t)
	response, err := indexer.Index("trumpet", "doc", "first", true, example_doc)

	// check for errors
	AssertNoError(t, "Failed to create doc", err)

	// was the response correct?
	expectedResponse := IndexResponse{"first", "trumpet", "doc", true}
	AssertEqual(t, expectedResponse == response, "Index response after create is incorrect", expectedResponse, response)

	// is the url valid?
	_, err = url.ParseRequestURI(httpMock.RequestURL)
	AssertNoError(t, "Illegal URL", err)

	// did it post where we wanted?
	expectedURL := indexer.docURL("trumpet", "doc", "first")
	AssertTrue(t, "URL does not contain docURL",
		strings.Contains(httpMock.RequestURL, expectedURL))

	// is it valid json?
	var requestDoc interface{}
	err = json.Unmarshal(httpMock.RequestBody, &requestDoc)
	AssertNoError(t, "request was invalid JSON!", err)

	// we won't check that the requestDoc is definitely correct because
	// there are many ways to send something to ES and this might make
	// our tests too brittle. round-trip integration tests can check that
}
