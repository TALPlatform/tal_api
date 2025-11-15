package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DbTracer struct {
	isDevelopment bool
}

func NewDbTracer(isDevelopment bool) *DbTracer {
	return &DbTracer{
		isDevelopment: isDevelopment,
	}
}

const maxArgSize = 256

func sanitizeArgs(args []any) []any {
	if len(args) == 0 {
		return args
	}

	sanitized := make([]any, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case []byte:
			if len(v) > maxArgSize {
				sanitized[i] = fmt.Sprintf("[...%d bytes truncated...]", len(v))
			} else {
				sanitized[i] = v
			}
		case string:
			if len(v) > maxArgSize {
				sanitized[i] = fmt.Sprintf("%s [...%d chars truncated...]", v[:maxArgSize], len(v)-maxArgSize)
			} else {
				sanitized[i] = v
			}
		default:
			sanitized[i] = arg
		}
	}
	return sanitized
}
func (tracer *DbTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.Info()
	if tracer.isDevelopment {
		sanitizedArgs := sanitizeArgs(data.Args)
		logger.Interface("arguments", sanitizedArgs).
			Str("query", data.SQL).
			Msg("DB Call Start")

	}
	return ctx
}

func (tracer *DbTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.Info()
	if tracer.isDevelopment {
		logger.Interface("arguments", data).
			Err(data.Err).
			Msg("DB Call End")
	}
}
