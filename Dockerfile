# syntax=docker/dockerfile:1

## Build
FROM golang:1.21.7-bookworm AS build

WORKDIR /app

COPY . .

COPY . .
RUN go mod tidy
RUN go build -o signet

## Deploy
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /app

COPY --from=build /app/signet signet

EXPOSE 8080

ENTRYPOINT ["./signet"]
