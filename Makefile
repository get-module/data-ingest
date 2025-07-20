APP_NAME=module-ingest

.PHONY: all build run clean test worker recover

all: build

build:
	go build -o bin/ingest cmd/ingest-server/main.go
	go build -o bin/worker cmd/queue-worker/main.go
	go build -o bin/recover cmd/recovery/main.go

run:
	go run cmd/ingest-server/main.go

worker:
	go run cmd/queue-worker/main.go

recover:
	go run cmd/recover-dump/main.go

test:
	go test ./...

clean:
	rm -rf bin/
