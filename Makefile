LOCAL_BIN:=$(CURDIR)/bin

format:
	golangci-lint cache clean
	golangci-lint run --fix ./...


build:
	go build -buildvcs=false -o $(LOCAL_BIN)/migrator ./cmd/migrator

run-up-dev:
	docker-compose -f docker/docker-compose.yml stop && docker-compose -f docker/docker-compose.yml up --build -d


MOCKS := mocks
MOCKGEN := github.com/vektra/mockery/v2@v2.42.1

gen_mocks: clean-mocks
	go run $(MOCKGEN) --dir ./internal --all --output $(MOCKS) --with-expecter --keeptree --case snake

clean-mocks:
	rm -rf $(MOCKS)