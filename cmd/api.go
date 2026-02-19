package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/psds-microservice/recording-service/internal/application"
	"github.com/psds-microservice/recording-service/internal/config"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run gRPC server (IngestStream)",
	RunE:  runAPI,
}

func runAPI(cmd *cobra.Command, args []string) error {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	app, err := application.New(cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return app.Run(ctx)
}
