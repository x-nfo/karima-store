# STAGE 1: Membangun aplikasi (Builder)
FROM golang:1.21-alpine AS builder
WORKDIR /app
# Install git & ca-certificates untuk keperluan fetch library dan HTTPS
RUN apk add --no-cache git ca-certificates
# Copy dependensi dulu agar caching layer Docker efisien
COPY go.mod go.sum ./
RUN go mod download
# Copy seluruh kode dan build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# STAGE 2: Menjalankan aplikasi (Final Image)
FROM alpine:latest AS final
WORKDIR /app
# Copy file binary dari stage builder
COPY --from=builder /app/main .
# Copy folder migrations (Penting untuk auto-migrate)
COPY --from=builder /app/migrations ./migrations
# Copy .env jika diperlukan (opsional, biasanya via env_file di compose)
COPY .env .
EXPOSE 8080
CMD ["./main"]
