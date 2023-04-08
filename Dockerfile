ARG GOVERSION="1.20.0"
FROM golang:${GOVERSION}-bullseye AS dev

WORKDIR /app

RUN apt-get update -qq && apt-get install -y curl
