LOCAL_BIN:=$(CURDIR)/bin

format:
	golangci-lint cache clean
	golangci-lint run --fix ./...


build:
	go build -o $(LOCAL_BIN)/migrator ./cmd/migrator

run-up-dev:
	docker-compose -p $(PROJECT_NAME) -f docker/docker-compose.yml up --build -d
