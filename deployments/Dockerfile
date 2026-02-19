FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /recording-service ./cmd/recording-service

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /recording-service .
# Default storage dir inside container; override with volume or env
ENV STORAGE_DIR=/app/recordings
EXPOSE 8096
CMD ["./recording-service"]
