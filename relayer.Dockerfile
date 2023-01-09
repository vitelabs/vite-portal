# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as build

WORKDIR /buildroot

ADD . .

WORKDIR /buildroot/relayer

RUN go mod download

RUN go build -o relayer cmd/relayer/main.go

FROM alpine:3.16

## Install bash shell
RUN apk update \
    && apk upgrade \
    && apk add --no-cache bash \
    bash-doc \
    bash-completion \
    gettext \
    && rm -rf /var/cache/apk/* \
    && /bin/bash

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=build /buildroot/relayer .

EXPOSE 56331 56332 56333 56334

ENTRYPOINT ["./relayer"]