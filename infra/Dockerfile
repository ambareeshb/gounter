FROM golang:1.18.3-alpine3.15 as builder

WORKDIR /root

RUN apk add --no-cache git
RUN apk add build-base

ENV GO111MODULE "on"

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY . /build
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app ./cmd

FROM alpine:3.15.0 as base
COPY --from=builder /build/app /gounter

EXPOSE 8081 
CMD ["./gounter"]
