package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func Init(env string) {
	level := slog.LevelInfo
	if env == "development" || env == "dev" {
		level = LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)
				if src != nil {
					short := src.File
					for i := len(short) - 1; i >= 0; i-- {
						if short[i] == '/' || short[i] == '\\' {
							short = short[i+1:]
							break
						}
					}
					a.Value = slog.StringValue(short)
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch env {
	case "production", "prod":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func Source() slog.Attr {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return slog.Attr{}
	}
	return slog.String("source", formatSource(file, line))
}

func Duration(d time.Duration) slog.Attr {
	return slog.Duration("duration", d)
}

func formatSource(file string, line int) string {
	short := file
	for i := len(short) - 1; i >= 0; i-- {
		if short[i] == '/' || short[i] == '\\' {
			short = short[i+1:]
			break
		}
	}
	return short + ":" + itoa(line)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
