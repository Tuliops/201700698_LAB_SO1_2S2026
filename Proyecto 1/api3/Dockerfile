#  Compilación
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

#  Imagen Final OCI
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8083
CMD ["./server"]
