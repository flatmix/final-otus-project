FROM golang:1.22-bullseye

COPY . /go/src/migrator

RUN apt update && \
    apt install libc6 && \
    cd /go/src/migrator && \
    make build

WORKDIR /go/src/migrator