FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY services/notification-service/go.* ./

COPY shared /shared

RUN go mod download

COPY services/notification-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o notification-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/notification-service .

COPY services/notification-service/migrations ./migrations

CMD ["./notification-service"]