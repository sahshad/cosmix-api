FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY services/communication-service/go.* ./

COPY shared /shared

RUN go mod download

COPY services/communication-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o communication-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/communication-service .

COPY services/communication-service/migrations ./migrations

CMD ["./communication-service"]
