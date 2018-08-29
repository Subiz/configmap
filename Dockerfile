FROM golang:1.11.0-alpine as builder
RUN apk update && apk add git
WORKDIR /go/src/git.subiz.net/configmap/
COPY ./ ./
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build
FROM alpine:3.8
COPY --from=builder /go/src/git.subiz.net/configmap/configmap /usr/local/bin/
