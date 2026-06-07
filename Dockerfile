# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/server

# ---- Run Stage ----
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /auth-service .

# Copy migrations so the app can run them on startup if needed
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./auth-service"]
