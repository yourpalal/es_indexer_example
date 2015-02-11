package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

// uses the following vars from other tests
// indexer *ElasticSearchIndexer (from elasticsearch_unit_test.go)
// example_doc Document (from elasticsearch_unit_test.go)

const (
	ES_TEST_INDEX = "testing"
)

// setup_es_integration_tests setus up a real indexer or skips a test if
// the environment variables are not set
func setup_es_integration_tests(t *testing.T) {
	indexer = MakeElasticSearchIndexerFromEnv(&http.Client{})
	if indexer == nil {
		t.Skip("Could not build ElasticSearchIndexer from env. variables (ES_HOSTNAME, ES_PORT")
	}
}

// ESDocResponse represents a very bare ElasticSearch response containing a
// single document
type ESDocResponse struct {
	Source *Document `json:"_source"`
}

// clear_es_data deletes all elastic search data from the given index (on the
// server referred to by indexer)
func clear_es_data(t *testing.T, index string) {
	index_url := fmt.Sprintf("http://%s:%s/%s", indexer.host, indexer.port, index)
	req, err := http.NewRequest("DELETE", index_url, strings.NewReader(""))
	AssertNoError(t, "Error when making DELETE request", err)

	client := &http.Client{}
	resp, err := client.Do(req)
	AssertNoError(t, "Error when clearing elastic search data", err)
	defer resp.Body.Close()
}

// es_get_doc is for testing purposes, it gets a doc from ES (as referred to by indexer)
func es_get_doc(t *testing.T, index string, _type string, id string) *Document {
	url := indexer.docURL(index, _type, id)
	response, err := http.Get(url)
	AssertNoError(t, "failed to get doc from Elastic Search", err)
	defer response.Body.Close()

	AssertEqual(t, response.StatusCode == 200, "HTTP error getting doc",
		200, response.StatusCode)

	json_buffer, err := ioutil.ReadAll(response.Body)
	AssertNoError(t, "failed to read response from Elastic Search", err)

	var esdoc *ESDocResponse
	err = json.Unmarshal(json_buffer, &esdoc)
	AssertNoError(t, "failed to parse response from Elasticsearch", err)
	if esdoc == nil || esdoc.Source == nil {
		t.Fatalf("insufficient json response from Elasticsearch %v", esdoc)
	}

	return esdoc.Source
}

// Test_MakeElasticSearchIndexer tests making an ElasticSearchIndexer instance
// based on environment variables
func Test_MakeElasticSearchIndexer(t *testing.T) {
	// capture original values, reset with defer
	es_host_env := os.Getenv(ES_HOSTNAME_KEY)
	defer os.Setenv(ES_HOSTNAME_KEY, es_host_env)
	es_port_env := os.Getenv(ES_PORT_KEY)
	defer os.Setenv(ES_PORT_KEY, es_port_env)

	// unset the environment variables
	os.Unsetenv(ES_HOSTNAME_KEY)
	os.Unsetenv(ES_PORT_KEY)

	indexer = MakeElasticSearchIndexerFromEnv(&http.Client{})
	AssertTrue(t, "made indexer without getting data from the environment?", indexer == nil)

	// create test values in the environment
	os.Setenv(ES_HOSTNAME_KEY, "test_hostname")
	os.Setenv(ES_PORT_KEY, "1546")
	indexer = MakeElasticSearchIndexerFromEnv(&http.Client{})
	AssertTrue(t, "MakeElasticSearchIndexerFromEnv failed to read env", indexer != nil)

	AssertEqual(t, indexer.host == "test_hostname", "wrong hostname",
		indexer.host, "test_hostname")
	AssertEqual(t, indexer.port == "1546", "wrong port", indexer.port, "1546")
}

// Test_CreateUpdateDoc tests creating, attempting to recreate, and
// updating a doc with the ElasticSearchIndexer, backed by a real
// Elasticsearch server.
func Test_CreateUpdateDoc(t *testing.T) {
	setup_es_integration_tests(t)
	clear_es_data(t, ES_TEST_INDEX)

	result, err := indexer.Index(ES_TEST_INDEX, "docs", "1", true, example_doc)
	AssertNoError(t, "error creating doc", err)
	expectedResponse := IndexResponse{"1", ES_TEST_INDEX, "docs", true}
	AssertEqual(t, result == expectedResponse, "incorrect response on create",
		expectedResponse, result)

	stored := es_get_doc(t, ES_TEST_INDEX, "docs", "1")
	AssertEqual(t, *stored == example_doc, "indexed doc is incorrect",
		example_doc, *stored)

	modified_doc := Document{
		Title:     example_doc.Title,
		Body:      example_doc.Body + " hope you like it!",
		Timestamp: example_doc.Timestamp,
	}

	// we should not be able to recreate a doc
	result, err = indexer.Index(ES_TEST_INDEX, "docs", "1", true, example_doc)
	AssertTrue(t, "Bad: succeeded in creating a doc twice", err != nil)
	stored = es_get_doc(t, ES_TEST_INDEX, "docs", "1")
	AssertEqual(t, *stored == example_doc, "failed create caused update!",
		example_doc, *stored)

	// but we should be able to update one!
	result, err = indexer.Index(ES_TEST_INDEX, "docs", "1", false, modified_doc)
	AssertNoError(t, "error updating doc", err)
	expectedResponse = IndexResponse{"1", ES_TEST_INDEX, "docs", false}
	AssertEqual(t, result == expectedResponse, "incorrect response on update",
		expectedResponse, result)

	stored = es_get_doc(t, ES_TEST_INDEX, "docs", "1")
	AssertEqual(t, *stored == modified_doc, "updated doc is incorrect",
		modified_doc, *stored)
}

// Test_UpdateNewDoc tests attempting to update a doc that has not yet
// been indexed with the ElasticSearchIndexer, backed by a real
// Elasticsearch server.
func Test_UpdateNewDoc(t *testing.T) {
	setup_es_integration_tests(t)
	clear_es_data(t, ES_TEST_INDEX)

	_, err := indexer.Index(ES_TEST_INDEX, "docs", "1", false, example_doc)
	AssertTrue(t, "Bad: succeeded in updating an unindexed doc", err != nil)
}
