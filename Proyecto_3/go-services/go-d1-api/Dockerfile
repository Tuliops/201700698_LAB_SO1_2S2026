FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]