.PHONY: deps build test

BINARY := govuk_crawler_worker

all: deps build test

deps:
	go run third_party.go get -t -v .

build:
	go run third_party.go build -o $(BINARY)

test:
	go run third_party.go test -v \
		github.com/alphagov/govuk_crawler_worker \
		github.com/alphagov/govuk_crawler_worker/http_crawler \
		github.com/alphagov/govuk_crawler_worker/queue \
		github.com/alphagov/govuk_crawler_worker/ttl_hash_set \
