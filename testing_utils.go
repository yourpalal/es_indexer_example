package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// a few mini utils for easier testing

// MockPoster implements HTTPPoster by saving the request data and returning
// a static response
type MockPoster struct {
	Result http.Response
	Err    error

	RequestURL      string
	RequestBody     []byte
	RequestBodyType string
}

// Post returns the Result and Err members while reading and saving the request
// data.
func (mock *MockPoster) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	mock.RequestURL = url
	mock.RequestBodyType = bodyType
	var readErr error
	mock.RequestBody, readErr = ioutil.ReadAll(body)
	if readErr != nil {
		log.Fatalf("failed to ReadAll in MockPoster: %s", readErr.Error())
	}

	return &mock.Result, mock.Err
}

// AssertTrue fails a test if its ok argument is false, and logs
// the provided message
func AssertTrue(t *testing.T, msg string, ok bool) {
	if !ok {
		t.Fatalf(msg)
	}
}

// AssertFalse fails a test if its ok argument is true , and logs
// the provided message
func AssertFalse(t *testing.T, msg string, ok bool) {
	if ok {
		t.Fatalf(msg)
	}
}

// AssertNoError will fail a test with the provided message and
// the error message if err is not nil
func AssertNoError(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatalf("%s\n error: %v", msg, err.Error())
	}
}

// Give a nice error message in the case of non-equality. We still make
// the caller do the comparison for us because that's just a bit easier
func AssertEqual(t *testing.T, equals bool, msg string, expected interface{}, actual interface{}) {
	if !equals {
		t.Fatalf("%s expected: %v got: %v", msg, expected, actual)
	}
}
