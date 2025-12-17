FROM golang:1.21-alpine AS builder


ENV CGO_ENABLED=0
ENV TZ=Europe/London

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o /goapp cmd/app/main.go



FROM alpine:3.18

 
RUN apk add --no-cache tzdata bash

WORKDIR /app


COPY --from=builder /goapp /app/main


EXPOSE 8080


CMD ["/app/main"]