FROM golang:1.24-alpine AS builder

# Установка пакетов
RUN apk add --no-cache \
    git \
    bash \
    protobuf \
    protobuf-dev \
    curl \
    unzip

# Установка 2 утилит
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest
WORKDIR /root/

# Копируем бинарь
COPY --from=builder /app/main .

# Копируем конфиги
COPY --from=builder /app/configs ./configs

EXPOSE 50051
CMD ["./main"]
