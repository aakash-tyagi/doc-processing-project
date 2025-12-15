# --- Stage 1: Build ---
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ingestion ./ingestion

# --- Stage 2: Run ---
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/ingestion .

EXPOSE 8080

CMD ["./ingestion"]