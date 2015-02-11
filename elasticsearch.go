package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Poster abstracts some basic functionality of http.Client
// so that we can do dependency injection and testing of http requests
type HTTPPoster interface {
	Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
}

// ElasticSearchIndexer implements Indexer by indexing documents in an
// Elasticsearch instance.
type ElasticSearchIndexer struct {
	host, port string
	client     HTTPPoster
}

// ElasticSearchUpdateRequest represents the json-formatted request that
// Elasticsearch requires for document updates
type ElasticSearchUpdateRequest struct {
	Doc interface{} `json:"doc"`
}

// verify that ElasticSearchIndexer implements Indexer by writing documents to
// an Elasticsearch (http://elasticsearch.org) server.
var _ Indexer = ElasticSearchIndexer{}

const (
	// environment variable to read for hostname
	ES_HOSTNAME_KEY = "ES_HOSTNAME"

	// environment variable to read for port
	ES_PORT_KEY = "ES_PORT"
)

// create an instance of ElasticSearchIndexer that gets its hostname and
// port from environment variables (ES_HOSTNAME and ES_PORT). If these
// are not present in the environment, it will return nil
func MakeElasticSearchIndexerFromEnv(client HTTPPoster) *ElasticSearchIndexer {
	es_host_env := os.Getenv(ES_HOSTNAME_KEY)
	es_port_env := os.Getenv(ES_PORT_KEY)

	if es_port_env == "" || es_host_env == "" {
		return nil
	}

	return &ElasticSearchIndexer{es_host_env, es_port_env, client}
}

// Index indexes the provided data in ElasticSearch
// the data will be indexed in /<index>/<_type>/<id>
//  data should be marshallable with json.
//
// When create = false, a previous version of the document to be indexed
// must already have been indexed. An update request will be issued to
// Elasticsearch, and a failure will result (noted in the error return value)
// in attempting to update an unindexed document.
//
// Similarly, if create=true and the document exists, an error will be returned,
// and the document will not be updated.
func (indexer ElasticSearchIndexer) Index(index string, _type string, id string, create bool, data interface{}) (response IndexResponse, err error) {
	response = IndexResponse{id, index, _type, false}

	if !create {
		data = &ElasticSearchUpdateRequest{data}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return response, err
	}

	docURL := indexer.docURL(index, _type, id)
	if create {
		docURL += "/_create"
	} else {
		docURL += "/_update"
	}

	httpR, err := indexer.client.Post(docURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	if httpR.StatusCode < 200 || 300 <= httpR.StatusCode {
		err = fmt.Errorf("Bad http response from elastic search : %v", httpR)
		return
	}

	if create && httpR.StatusCode != 201 {
		err = errors.New("Index() Attempted create but did not get 201 status")
		return
	}

	return IndexResponse{id, index, _type, create}, nil
}

// docURL returns the elastic search url for a given document
// i.e you could CURL this to get the document
func (indexer *ElasticSearchIndexer) docURL(index string, _type string, id string) string {
	return fmt.Sprintf("http://%s:%s/%s/%s/%s", indexer.host, indexer.port, index, _type, id)
}
