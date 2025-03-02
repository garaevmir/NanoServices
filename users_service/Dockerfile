# Сборочный этап
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o users_service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/users_service .

EXPOSE 8081

CMD ["./users_service"]