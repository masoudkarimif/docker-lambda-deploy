FROM golang:1.19-alpine AS builder

LABEL maintainer="github.com/masoudkarimif"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o main cmd/*.go

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder ./build/main /

ENTRYPOINT ["/main"]