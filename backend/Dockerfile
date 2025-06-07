# Stage 1: build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main .

# Stage 2: create a clean image
FROM debian:bookworm-slim

# ✅ Install root certificates for HTTPS (required for Firebase)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY keys /app/keys
COPY .env ./

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["/app/main"]
