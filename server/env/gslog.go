package env

import (
	"context"
	"log/slog"
	"os"
)

type CloudLoggingHandler struct {
	handler slog.Handler
}

func slogToCloudLogging(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	} else if a.Key == slog.SourceKey {
		a.Key = "logging.googleapis.com/sourceLocation"
	} else if a.Key == slog.LevelKey {
		a.Key = "severity"
		level, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}
		switch level {
		case slog.LevelDebug:
			return slog.String("severity", "DEBUG")
		case slog.LevelInfo:
			return slog.String("severity", "INFO")
		case slog.LevelWarn:
			return slog.String("severity", "WARNING")
		case slog.LevelError:
			return slog.String("severity", "ERROR")
		default:
			return slog.String("severity", "DEFAULT")
		}
	}
	return a
}

func NewCloudLoggingHandler() *CloudLoggingHandler {
	return &CloudLoggingHandler{handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: slogToCloudLogging,
	})}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	// trace := traceFromContext(ctx)
	// if trace != "" {
	// 	rec = rec.Clone()
	// 	// Add trace ID	to the record so it is correlated with the Cloud Run request log
	// 	// See https://cloud.google.com/trace/docs/trace-log-integration
	// 	rec.Add("logging.googleapis.com/trace", slog.StringValue(trace))
	// }
	return h.handler.Handle(ctx, rec)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}

func SetupCloudLogging() {
	handler := slog.New(NewCloudLoggingHandler())
	slog.SetDefault(handler)
}
