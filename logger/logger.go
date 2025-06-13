package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var logger *slog.Logger

func init() {
	if err := initLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func initLogger() error {
	// Create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Get log level
	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create file writer
	logFile, err := os.OpenFile(
		filepath.Join("logs", "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return err
	}

	// Multi-writer (console + file)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Choose handler
	var handler slog.Handler
	if os.Getenv("ENVIRONMENT") == "production" {
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	} else {
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

func GetLogger() *slog.Logger {
	return logger
}

// ✅ ADD THESE MISSING HELPER FUNCTIONS:

// Helper functions without context
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// Helper functions with context
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}
