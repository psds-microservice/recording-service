# No local deps; build from repo root for consistency.
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY recording-service ./recording-service
WORKDIR /app/recording-service
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /recording-service ./cmd/recording-service

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /recording-service .
ENV STORAGE_DIR=/app/recordings
EXPOSE 8096
CMD ["./recording-service"]
