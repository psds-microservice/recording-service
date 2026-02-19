# recording-service

**EN:** Receives a copy of the video stream from streaming-service (gRPC client stream), writes it to local storage, and returns a playback URL. Streaming-service then sets this URL on the consultation session via session-manager `SetRecordingUrl`.  
**RU:** Принимает копию видеопотока от streaming-service (gRPC client stream), записывает в локальное хранилище и возвращает URL для воспроизведения. Streaming-service затем проставляет этот URL в сессии консультации через session-manager `SetRecordingUrl`.

## Structure (aligned with other PSDS services)

**EN:** Cobra CLI (`cmd/root.go`, `cmd/api.go`), `make build` / `make run` / `make run-dev`, `.env.example`, `deployments/Dockerfile` and `docker-compose.yml`, `.gitignore`.  
**RU:** Cobra CLI (`cmd/root.go`, `cmd/api.go`), `make build` / `make run` / `make run-dev`, `.env.example`, `deployments/Dockerfile` и `docker-compose.yml`, `.gitignore`.

## Run

**EN:** From `recording-service`: `make run` or `make run-dev`, or `go run ./cmd/recording-service api`. Copy `.env.example` to `.env` if needed.  
**RU:** Из каталога `recording-service`: `make run` или `make run-dev`, либо `go run ./cmd/recording-service api`. При необходимости скопируйте `.env.example` в `.env`.

```bash
# .env (see .env.example)
GRPC_PORT=8096
STORAGE_DIR=./recordings
RECORDING_BASE_URL=http://localhost:8096/recordings
```

## Integration

**EN:** streaming-service connects to recording-service when `ENABLE_RECORDING=true`, `RECORDING_SERVICE_ADDR=localhost:8096`, and `SESSION_MANAGER_GRPC_ADDR=localhost:9091`. It sends each client WebSocket chunk via gRPC `IngestStream` and on session end calls session-manager `SetRecordingUrl(stream_session_id, url)`.  
**RU:** streaming-service подключается к recording-service при `ENABLE_RECORDING=true`, `RECORDING_SERVICE_ADDR=localhost:8096` и `SESSION_MANAGER_GRPC_ADDR=localhost:9091`. Каждый чанк от клиента по WebSocket отправляется по gRPC `IngestStream`; по завершении сессии вызывается session-manager `SetRecordingUrl(stream_session_id, url)`.
