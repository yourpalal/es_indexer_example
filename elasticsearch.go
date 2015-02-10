package main

import (
	"encoding/json"
	"fmt"
)

// ElasticSearchIndexer implements Indexer by indexing documents in an
// Elasticsearch instance.
type ElasticSearchIndexer struct {
	host, port string
}

// verify that ElasticSearchIndexer implements Indexer by writing documents to
// an Elasticsearch (http://elasticsearch.org) server.
var _ Indexer = ElasticSearchIndexer{}

// Index indexes the provided data in ElasticSearch
// the data will be indexed in /<index>/<_type>/<id>
//  data should be marshallable with json
func (indexer ElasticSearchIndexer) Index(index string, _type string, id string, create bool, data interface{}) (IndexResponse, error) {
	_, err := json.Marshal(data)
	if err != nil {
		return IndexResponse{id, index, _type, false}, err
	}

	// TODO when creating, do op_type=create, set timestamp?
	return IndexResponse{id, index, _type, create}, nil
}

// docURL returns the elastic search url for a given document
// i.e you could CURL this to get the document
func (indexer *ElasticSearchIndexer) docURL(index string, _type string, id string) string {
	return fmt.Sprintf("http://%s:%s/%s/%s/%s", indexer.host, indexer.port, index, _type, id)
}
