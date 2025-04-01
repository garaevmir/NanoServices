FROM golang:1.23 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY proto/ ./proto/

RUN protoc --go_out=. --go_opt=module=github.com/nanoservices/gateway \
    --go-grpc_out=. --go-grpc_opt=module=github.com/nanoservices/gateway \
    proto/*.proto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]
