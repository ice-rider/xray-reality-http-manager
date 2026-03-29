FROM golang:1.25-alpine AS builder

WORKDIR /app

# Установка зависимостей для сборки с CGO (sqlite)
RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка приложения
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
