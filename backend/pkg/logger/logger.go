package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

type Config struct {
	Level       string
	ServiceName string
	CellID      string
	Pretty      bool
}

func New(cfg Config) *Logger {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var output = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("cell_id", cfg.CellID).
		Caller().
		Logger()

	return &Logger{logger: logger}
}

func (l *Logger) Debug(msg string)                                  { l.logger.Debug().Msg(msg) }
func (l *Logger) Debugf(format string, args ...interface{})         { l.logger.Debug().Msgf(format, args...) }
func (l *Logger) Info(msg string)                                   { l.logger.Info().Msg(msg) }
func (l *Logger) Infof(format string, args ...interface{})          { l.logger.Info().Msgf(format, args...) }
func (l *Logger) Warn(msg string)                                   { l.logger.Warn().Msg(msg) }
func (l *Logger) Warnf(format string, args ...interface{})          { l.logger.Warn().Msgf(format, args...) }
func (l *Logger) Error(err error, msg string)                       { l.logger.Error().Err(err).Msg(msg) }
func (l *Logger) Errorf(err error, format string, args ...interface{}) { l.logger.Error().Err(err).Msgf(format, args...) }
func (l *Logger) Fatal(err error, msg string)                       { l.logger.Fatal().Err(err).Msg(msg) }

func (l *Logger) With(key string, value interface{}) *Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &Logger{logger: newLogger}
}
