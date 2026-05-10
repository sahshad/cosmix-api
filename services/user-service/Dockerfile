# Build stage
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY services/user-service/go.* ./

# Copy cosmix-events for messaging
COPY shared /shared

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY services/user-service/ .

# Copy .env
COPY services/user-service/.env .env

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o user-service ./cmd/main.go

# Runtime stage
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Install only ca-certificates
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/user-service .

# Copy migrations
COPY services/user-service/migrations ./migrations

# Run the binary
CMD ["./user-service"]
