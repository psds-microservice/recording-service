package application

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/psds-microservice/recording-service/internal/config"
	grpcimpl "github.com/psds-microservice/recording-service/internal/grpc"
	"github.com/psds-microservice/recording-service/pkg/gen/recording_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App is the recording-service gRPC application.
type App struct {
	cfg *config.Config
	srv *grpc.Server
	lis net.Listener
}

// New creates the application: validates config, creates gRPC server and listener.
func New(cfg *config.Config) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	logger, _ := zap.NewProduction()
	if cfg.AppEnv == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	lis, err := net.Listen("tcp", cfg.GRPCAddr())
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	recSrv := grpcimpl.NewServer(cfg, logger)
	recording_service.RegisterRecordingServiceServer(grpcServer, recSrv)
	reflection.Register(grpcServer)

	return &App{cfg: cfg, srv: grpcServer, lis: lis}, nil
}

// Run starts the gRPC server and blocks until ctx is cancelled.
func (a *App) Run(ctx context.Context) error {
	log.Printf("recording-service gRPC listening on %s", a.lis.Addr())
	go func() {
		if err := a.srv.Serve(a.lis); err != nil {
			log.Printf("grpc serve: %v", err)
		}
	}()
	<-ctx.Done()
	stopped := make(chan struct{})
	go func() {
		a.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		return nil
	case <-time.After(10 * time.Second):
		a.srv.Stop()
		return fmt.Errorf("grpc shutdown timeout")
	}
}
