LOCAL_BIN:=$(CURDIR)/bin

format:
	golangci-lint cache clean
	golangci-lint run --fix ./...


build:
	go build -buildvcs=false -o $(LOCAL_BIN)/migrator ./cmd/migrator

run-up-dev:
	docker-compose -f docker/docker-compose.yml stop && docker-compose -f docker/docker-compose.yml up --build -d
