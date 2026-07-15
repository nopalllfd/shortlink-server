FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Build final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin && \
    chmod +x /usr/local/bin/migrate

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
