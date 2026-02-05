FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
ENV GOTOOLCHAIN=auto
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/app/main.go


RUN CGO_ENABLED=0 GOOS=linux go build -o app-bin ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o notifier-bin ./cmd/notifier/main.go

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/app-bin ./main
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
EXPOSE 8000
CMD ["./main"]


FROM alpine:latest AS notifier
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/notifier-bin ./notifier
COPY --from=builder /app/.env .
CMD ["./notifier"]