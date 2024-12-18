package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

type LogLevel slog.Level

const (
	OffLevel     = LogLevel(slog.Level(1000))
	FatalLevel   = LogLevel(slog.Level(9))
	ErrorLevel   = LogLevel(slog.LevelError)
	WarnLevel    = LogLevel(slog.LevelWarn)
	InfoLevel    = LogLevel(slog.LevelInfo)
	DebugLevel   = LogLevel(slog.LevelDebug)
	HttpLevel    = LogLevel(slog.Level(-7))
	MetricsLevel = LogLevel(slog.Level(-8))
	TraceLevel   = LogLevel(slog.Level(-10))
)

func (l *LogLevel) String() string {
	switch *l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case FatalLevel:
		return "FATAL"
	case ErrorLevel:
		return "ERROR"
	case TraceLevel:
		return "TRACE"
	case WarnLevel:
		return "WARN"
	case MetricsLevel:
		return "METRICS"
	case HttpLevel:
		return "HTTP"
	case OffLevel:
		return "OFF"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int64(*l))
	}
}

func (l *LogLevel) Set(val string) error {
	switch strings.ToLower(val) {
	case "fatal":
		*l = FatalLevel
	case "error":
		*l = ErrorLevel
	case "warn":
		*l = WarnLevel
	case "info":
		*l = InfoLevel
	case "debug":
		*l = DebugLevel
	case "http":
		*l = HttpLevel
	case "metrics":
		*l = MetricsLevel
	case "trace":
		*l = TraceLevel
	case "off":
		*l = OffLevel
	default:
		return fmt.Errorf(fmt.Sprintf("unknown log level: %s", val))
	}
	return nil
}

func (l *LogLevel) Type() string {
	return "LogLevel"
}

func LogLevelPtr(logLevel LogLevel) *LogLevel {
	return &logLevel
}

func ConfigureLogs(machineReadable bool, logLevel LogLevel) {
	attrFunc := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			lvl := LogLevel(a.Value.Any().(slog.Level))
			a.Value = slog.StringValue(lvl.String())
		}
		return a
	}
	var handler slog.Handler = tint.NewHandler(os.Stdout, &tint.Options{
		Level:       slog.Level(logLevel),
		AddSource:   logLevel < InfoLevel,
		ReplaceAttr: attrFunc,
	})
	if machineReadable {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.Level(logLevel),
			AddSource:   logLevel < InfoLevel,
			ReplaceAttr: attrFunc,
		})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Debug(fmt.Sprintf("Logging started with level %s", logLevel.String()))

}

func Trace(msg string, args ...any) {
	slog.Log(context.Background(), slog.Level(TraceLevel), msg, args...)
}

func Metrics(msg string, args ...any) {
	slog.Log(context.Background(), slog.Level(MetricsLevel), msg, args...)
}

func Fatal(msg string, err error, args ...any) {
	things := append([]any{slog.String("error", err.Error())}, args...)
	slog.Log(context.Background(), slog.Level(FatalLevel), msg, things...)
}

type HttpData struct {
	Start      bool
	Method     string
	Url        string
	StatusCode int64
	Latencyms  int64
}

func Http(data HttpData, args ...any) {
	if data.Start {
		slog.Log(
			context.Background(),
			slog.Level(HttpLevel),
			fmt.Sprintf("Request %s:%s", data.Method, data.Url),
			args...)

	} else {
		slog.Log(
			context.Background(),
			slog.Level(HttpLevel),
			fmt.Sprintf("Response %s:%s", data.Method, data.Url),
			append(
				args,
				slog.Int64("status_code", data.StatusCode),
				slog.Int64("latency_ms", data.Latencyms),
			)...,
		)
	}
}
