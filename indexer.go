package main

// IndexResponse contains the response from the backend api regarding the status of the indexing
// operation. ID contains the reference to the specific document in the index, Index is the
// specific search index that the document exists on, Type is the specific type that the document
// is, Created is true when the document was created sucesfully, will be false if it was simply
// updated
type IndexResponse struct {
	ID      string
	Index   string
	Type    string
	Created bool
}

// Indexer provides an interface to the indexing service
type Indexer interface {
	// Index performs an indexing operation on data where it's type is _type, on the index
	// index, with the id being required. If the create flag is true, then the document will
	// be created in the index, if false it will update the document with data.
	Index(index string, _type string, id string, create bool, data interface{}) (IndexResponse, error)
}
