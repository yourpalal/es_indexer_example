package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// verify that ElasticSearchIndexer implements Indexer by writing documents to
// an Elasticsearch (http://elasticsearch.org) server.
var _ Indexer = ElasticSearchIndexer{}

const (
	ES_HOSTNAME_KEY = "ES_HOSTNAME"
	ES_PORT_KEY     = "ES_PORT"
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
//  data should be marshallable with json
func (indexer ElasticSearchIndexer) Index(index string, _type string, id string, create bool, data interface{}) (response IndexResponse, err error) {
	response = IndexResponse{id, index, _type, false}

	jsonData, err := json.Marshal(data)
	_, err = json.Marshal(data)
	if err != nil {
		return response, err
	}

	docURL := indexer.docURL(index, _type, id)
	if create {
		v := url.Values{}
		v.Set("op_type", "create")
		docURL += "?" + v.Encode()
	}

	httpR, err := indexer.client.Post(docURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
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
