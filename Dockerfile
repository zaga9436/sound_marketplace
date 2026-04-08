FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /soundmarket ./cmd/api

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /soundmarket /usr/local/bin/soundmarket
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["soundmarket"]
