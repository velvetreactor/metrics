FROM golang:1.21.3-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o metrics_server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/metrics_server /metrics_server
CMD ["/metrics_server"]