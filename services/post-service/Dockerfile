FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY services/post-service/go.* ./

COPY shared /shared

RUN go mod download

COPY services/post-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o post-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/post-service .

COPY services/post-service/migrations ./migrations

CMD ["./post-service"]