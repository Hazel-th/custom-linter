package ok

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func good(logger *slog.Logger, z *zap.Logger, sugar *zap.SugaredLogger) {
	slog.Info("starting server on port 8080")
	logger.ErrorContext(context.Background(), "failed to connect to database")
	z.Warn("connection failed")
	sugar.Infow("request completed", "request_id", "123")
	slog.Info("token validated")
}
