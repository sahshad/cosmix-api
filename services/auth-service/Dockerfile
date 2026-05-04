# Build stage
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY services/auth-service/go.* ./

# Copy cosmix-events for messaging
COPY cosmix-events /cosmix-events

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY services/auth-service/ .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o auth-service ./cmd/main.go

# Runtime stage
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Install only ca-certificates
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/auth-service .

# Copy migrations
COPY services/auth-service/migrations ./migrations

# Run the binary
CMD ["./auth-service"]
