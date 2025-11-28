FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-app ./cmd/mini-wallet

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/wallet-app /app/wallet-app
COPY --from=builder /app/config /app/config
COPY --from=builder /app/internal/infrastructure/storage/postgres/migrations \
    /app/internal/infrastructure/storage/postgres/migrations

EXPOSE 8080
CMD ["./wallet-app"]