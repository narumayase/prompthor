# Etapa 1: build
FROM golang:1.24-bullseye AS builder

# Librerías necesarias para confluent-kafka-go
RUN apt-get update && apt-get install -y \
    git \
    pkg-config \
    librdkafka-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o anyompt .

# Etapa 2: imagen final
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/anyompt .

# Librerías runtime + certificados raíz
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 8100
CMD ["./anyompt"]