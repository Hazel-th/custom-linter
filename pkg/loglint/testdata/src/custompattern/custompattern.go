package custompattern

import "log/slog"

func customRegexPattern() {
	slog.Info("order 1234 processed") // want "may contain sensitive data"
}
