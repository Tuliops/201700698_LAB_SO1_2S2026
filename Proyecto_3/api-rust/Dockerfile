FROM rust:1.90-slim AS builder
WORKDIR /app
COPY Cargo.toml .
COPY src ./src
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/target/release/api-rust .
EXPOSE 8080
CMD ["./api-rust"]
