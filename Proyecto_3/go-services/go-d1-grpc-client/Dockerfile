FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o client main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/client .
CMD ["./client"]