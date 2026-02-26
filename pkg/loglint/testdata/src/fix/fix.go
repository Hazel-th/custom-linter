package fix

import "log/slog"

func bad(logger *slog.Logger, token string, password string) {
	logger.Info("Starting server")            // want "start with a lowercase letter"
	logger.Info("ошибка подключения")         // want "contain only English language"
	logger.Info("connection failed!!!")       // want "must not contain special symbols or emoji"
	logger.Info("token: " + token)            // want "may contain sensitive data"
	logger.Info("user password: " + password) // want "may contain sensitive data"
}
