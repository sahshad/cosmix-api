FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY services/realtime-service/go.* ./

COPY shared /shared

RUN go mod download

COPY services/realtime-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o realtime-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/realtime-service .

EXPOSE 50055

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -qO- http://localhost:50055/health || exit 1

CMD ["./realtime-service"]
