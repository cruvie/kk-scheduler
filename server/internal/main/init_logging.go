package main

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		configSlog.Logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func unaryLogging() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(interceptorLogger())
}
