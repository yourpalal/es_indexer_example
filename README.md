ElasticSearch indexer implementation.
==========================================

Writes documents to an ElasticSearch server. This project was written and tested
in Go 1.4. To run unit tests, use

    go test

to run unit and integration tests, you must provide an elasticsearch server.
The testing data will be creating in the index "testing" and type "docs". You
can run both kinds of tests using

    ES_HOSTNAME=<your elasticsearch host> ES_PORT=<your elasticsearch port> go test

For instance, to run on a development machine with default settings, you might run:

    ES_HOSTNAME=localhost ES_PORT=9200 go test

If these environment variables are not present, the integration tests will be
skipped.


Notes on Implementation
-----------------------

I would generally use popular open-source libraries for a task like
this (one for nicer testing, maybe a mocking library, and one for
talking to Elasticsearch). However, in this case I decided not to as
this makes the resulting code more indicative of what I can accomplish
in Go.
