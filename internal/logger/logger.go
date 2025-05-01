package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fungicibus/order/config"
	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(cfg *config.Config, nonConsoleWriter io.Writer) (*Logger, error) {
	var writer io.Writer = os.Stdout
	if nonConsoleWriter != nil {
		writer = zerolog.MultiLevelWriter(os.Stdout, nonConsoleWriter)
	}

	zerolog.MessageFieldName = "_msg"
	stream := fmt.Sprintf("app=%s,env=%s", cfg.App.Name, cfg.App.Env)
	logger := zerolog.New(writer).With().
		Str("_stream", stream).
		Timestamp().
		Logger().
		Level(zerolog.Level(cfg.Log.Level))
	return &Logger{logger}, nil
}

func (l *Logger) SetLevel(level int) {
	l.Logger = l.Logger.Level(zerolog.Level(level))
}

func WtihSource(initialLogger *Logger, source string) *Logger {
	loggerWithSource := initialLogger.With().Str("from", source).Logger()
	return &Logger{
		Logger: loggerWithSource,
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v)
}
