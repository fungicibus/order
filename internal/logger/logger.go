package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/fungicibus/order/config"
)

type Logger struct {
	zerolog.Logger
}

func New(cfg *config.Config, nonConsoleWriters ...io.Writer) (*Logger, error) {
	var writer io.Writer = os.Stdout
	if len(nonConsoleWriters) > 0 {
		writers := make([]io.Writer, 0, len(nonConsoleWriters)+1)
		writers = append(writers, os.Stdout)
		writers = append(writers, nonConsoleWriters...)
		writer = zerolog.MultiLevelWriter(writers...)
	}

	stream := fmt.Sprintf("app=%s,env=%s", cfg.App.Name, cfg.App.Env)
	logger := zerolog.New(writer).With().
		Str("_stream", stream).
		Timestamp().
		Logger().
		Level(zerolog.Level(cfg.Log.Level))
	return &Logger{logger}, nil
}

func (l *Logger) SetLevel(level int) {
	l.Logger = l.Level(zerolog.Level(level))
}

func WithSource(initialLogger *Logger, source string) *Logger {
	loggerWithSource := initialLogger.With().Str("from", source).Logger()
	return &Logger{
		Logger: loggerWithSource,
	}
}

// Implements goose.Logger
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal().Msgf(format, v...)
}
