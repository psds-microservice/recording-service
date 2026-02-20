package grpc

import (
	"io"
	"sync"

	"github.com/psds-microservice/recording-service/internal/config"
	"github.com/psds-microservice/recording-service/internal/storage"
	"github.com/psds-microservice/recording-service/pkg/gen/recording_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements recording_service.RecordingServiceServer.
type Server struct {
	recording_service.UnimplementedRecordingServiceServer
	cfg *config.Config
	log *zap.Logger
	mu  sync.Mutex
	// active writers by session_id (one IngestStream per session at a time)
	writers map[string]*storage.SessionWriter
}

// NewServer creates the gRPC server.
func NewServer(cfg *config.Config, log *zap.Logger) *Server {
	return &Server{cfg: cfg, log: log, writers: make(map[string]*storage.SessionWriter)}
}

// IngestStream receives chunks from the client, writes to a file, returns URL when stream ends.
func (s *Server) IngestStream(stream recording_service.RecordingService_IngestStreamServer) error {
	var sessionID string
	var writer *storage.SessionWriter
	defer func() {
		if writer != nil {
			_ = writer.Close()
			s.mu.Lock()
			delete(s.writers, sessionID)
			s.mu.Unlock()
		}
	}()

	for {
		// Check if stream context is cancelled (client disconnected)
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
		}
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if chunk.SessionId == "" {
			return status.Error(codes.InvalidArgument, "session_id required")
		}
		if sessionID == "" {
			sessionID = chunk.SessionId
			path := s.cfg.RecordingPath(sessionID)
			w, err := storage.NewSessionWriter(path, s.log)
			if err != nil {
				return status.Errorf(codes.Internal, "create recording file: %v", err)
			}
			writer = w
			s.mu.Lock()
			s.writers[sessionID] = writer
			s.mu.Unlock()
		}
		if len(chunk.Data) > 0 {
			if _, err := writer.Write(chunk.Data); err != nil {
				return status.Errorf(codes.Internal, "write chunk: %v", err)
			}
		}
		if chunk.Last {
			break
		}
	}

	if writer == nil {
		return status.Error(codes.InvalidArgument, "no chunks received")
	}
	if err := writer.Close(); err != nil {
		s.log.Warn("close recording file", zap.Error(err))
	}
	writer = nil
	s.mu.Lock()
	delete(s.writers, sessionID)
	s.mu.Unlock()

	url := s.cfg.RecordingURL(sessionID)
	return stream.SendAndClose(&recording_service.RecordingResult{RecordingUrl: url})
}
