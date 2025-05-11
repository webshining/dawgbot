FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/bin/telegram ./cmd/telegram/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/bin/discord ./cmd/discord/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/telegram /app/
COPY --from=builder /app/bin/discord /app/

COPY .env /app/.env