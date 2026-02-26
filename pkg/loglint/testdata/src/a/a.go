package a

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func bad(logger *slog.Logger, z *zap.Logger, sugar *zap.SugaredLogger, password, apiKey, token string) {
	slog.Info("Starting server on port 8080")                         // want "start with a lowercase letter"
	logger.Error("ошибка подключения к базе данных")                  // want "contain only English language"
	z.Warn("connection failed!!!")                                    // want "must not contain special symbols or emoji"
	slog.WarnContext(context.Background(), "something went wrong...") // want "must not contain special symbols or emoji"
	sugar.Infow("token: " + token)                                    // want "may contain sensitive data"
	slog.Info("api_key=" + apiKey)                                    // want "may contain sensitive data"
	slog.Info("user password: " + password)                           // want "may contain sensitive data"
	z.Info("server started 🚀")                                        // want "must not contain special symbols or emoji"
	zap.L().Error("Failed to connect to database")                    // want "start with a lowercase letter"
}
