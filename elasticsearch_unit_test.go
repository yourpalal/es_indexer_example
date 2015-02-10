package main

import (
	"testing"
)


// unjsonable is a map[int]int that cannot be jsonified
var (
    unjsonable map[int]int
    indexer ElasticSearchIndexer
)

// setup to run before each test
func setup(t *testing.T) {

    unjsonable = make(map[int]int)
    unjsonable[1] = 1

    indexer = ElasticSearchIndexer{"127.0.0.1", "9200"}
}

func Test_ESIndexerReturnsJsonErrors(t *testing.T) {
    setup(t)

    r, err := indexer.Index("testing", "impossible", "0", true, unjsonable)

    if err == nil {
		t.Fatal("Failed to return json error")
    }
    if r.Created {
		t.Fatal("Set Created to true on failure")
    }
}

func Test_ESIndexer_docURL(t *testing.T) {
    setup(t)

    indexer.ip = "127.0.0.1"
    indexer.port = "9200"
    url := indexer.docURL("twitter", "tweet", "best_tweet");
    expected := "http://127.0.0.1:9200/twitter/tweet/best_tweet"
    if url != expected {
        t.Fatalf("Failed to build document URL. expected %s, got %s", expected, url)
    }
}
