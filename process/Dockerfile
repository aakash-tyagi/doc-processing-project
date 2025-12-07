# --- Stage 1: Build ---
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o process .

# --- Stage 2: Run ---
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/process .

EXPOSE 8081

CMD [ "./process" ]
