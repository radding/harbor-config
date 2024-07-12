package telemetry

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogLevel zerolog.Level

func (l *LogLevel) String() string {
	switch *l {
	case LogLevel(zerolog.DebugLevel):
		return "debug"
	case LogLevel(zerolog.InfoLevel):
		return "info"
	case LogLevel(zerolog.PanicLevel):
		return "panic"
	case LogLevel(zerolog.FatalLevel):
		return "fatal"
	case LogLevel(zerolog.ErrorLevel):
		return "error"
	case LogLevel(zerolog.TraceLevel):
		return "trace"
	case LogLevel(zerolog.WarnLevel):
		return "warn"
	default:
		return "unknown"
	}
}

func (l *LogLevel) Set(val string) error {
	switch strings.ToLower(val) {
	case "info":
		*l = LogLevel(zerolog.InfoLevel)
	case "warn":
		*l = LogLevel(zerolog.WarnLevel)
	case "panic":
		*l = LogLevel(zerolog.PanicLevel)
	case "fatal":
		*l = LogLevel(zerolog.FatalLevel)
	case "error":
		*l = LogLevel(zerolog.ErrorLevel)
	case "debug":
		*l = LogLevel(zerolog.DebugLevel)
	case "trace":
		*l = LogLevel(zerolog.TraceLevel)
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
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))
	if !machineReadable {
		out := zerolog.ConsoleWriter{Out: os.Stdout}
		out.PartsOrder = []string{
			"Identifier",
			"time",
			"level",
			"message",
		}
		out.FieldsExclude = []string{
			"Identifier",
		}
		out.FormatFieldValue = func(i interface{}) string {
			if i == nil {
				return ""
			}
			return fmt.Sprintf("%s", i)
		}

		log.Logger = log.Output(out)
	}
	log.Trace().Msgf("Starting logging with level: %s", logLevel.String())

}
