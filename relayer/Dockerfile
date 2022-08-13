# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as build

WORKDIR /buildroot

COPY go.mod ./
COPY go.sum ./

RUN go mod download

ADD . .

RUN go build -o relayer cmd/relayer/main.go

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

COPY --from=build /buildroot/relayer .

EXPOSE 56331 56332

ENTRYPOINT ["./relayer"]