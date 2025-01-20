# Build stage
FROM golang:1.23.4-alpine as builder
WORKDIR /app
COPY . .
COPY .env /app/.env
RUN go build -o /app/app

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app /app
COPY --from=builder /app/.env /app/.env
CMD ["/app/app"]
