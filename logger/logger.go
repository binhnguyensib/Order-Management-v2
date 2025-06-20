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
		// ✅ Không panic, fallback to stdout
		setupFallbackLogger()
	}
}

func initLogger() error {
	// ✅ Check if running in serverless environment
	if isServerlessEnvironment() {
		setupServerlessLogger()
		return nil
	}

	// ✅ Local development - try to create logs directory
	return setupLocalLogger()
}

func isServerlessEnvironment() bool {
	return os.Getenv("VERCEL") == "1" ||
		os.Getenv("LAMBDA_TASK_ROOT") != "" ||
		os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func setupServerlessLogger() {
	// ✅ Serverless: chỉ log ra stdout
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

	// ✅ JSON format cho production serverless
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false, // Disable source for cleaner serverless logs
	})

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func setupLocalLogger() error {
	// ✅ Local development - tạo logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		// Fallback to stdout nếu không tạo được
		setupFallbackLogger()
		return nil // Không return error, chỉ fallback
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
		// Fallback to stdout nếu không tạo được file
		setupFallbackLogger()
		return nil
	}

	// Multi-writer (console + file) cho local
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

func setupFallbackLogger() {
	// ✅ Fallback logger khi không thể setup file logging
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func GetLogger() *slog.Logger {
	if logger == nil {
		// ✅ Safety check - tạo fallback nếu logger nil
		setupFallbackLogger()
	}
	return logger
}

// ✅ Helper functions without context
func Debug(msg string, args ...any) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if logger != nil {
		logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if logger != nil {
		logger.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if logger != nil {
		logger.Error(msg, args...)
	}
}

// ✅ Helper functions with context
func DebugContext(ctx context.Context, msg string, args ...any) {
	if logger != nil {
		logger.DebugContext(ctx, msg, args...)
	}
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	if logger != nil {
		logger.InfoContext(ctx, msg, args...)
	}
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	if logger != nil {
		logger.WarnContext(ctx, msg, args...)
	}
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	if logger != nil {
		logger.ErrorContext(ctx, msg, args...)
	}
}
