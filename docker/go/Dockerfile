FROM golang:1.22-bullseye

COPY . /go/src/migrator

RUN apt update && \
	go install github.com/boumenot/gocover-cobertura@v1.2.0


WORKDIR /go/src/migrator