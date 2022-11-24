# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as build

WORKDIR /buildroot

ADD . .

WORKDIR /buildroot/orchestrator

RUN go mod download

RUN go build -o orchestrator cmd/orchestrator/main.go

FROM alpine:3.16

## Install bash shell
RUN apk update \
    && apk upgrade \
    && apk add --no-cache bash \
    bash-doc \
    bash-completion \
    && rm -rf /var/cache/apk/* \
    && /bin/bash

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=build /buildroot/orchestrator .

EXPOSE 57331 57332

ENTRYPOINT ["./orchestrator"]