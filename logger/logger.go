package logger

// TODO:
// take trace context as input to dump relevant info in error logs
import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AppLogger struct {
	logger zerolog.Logger
}

var appLogger *AppLogger = nil

func NewLogger() *AppLogger {
	// ensures that only one logger is created across the system
	if appLogger != nil {
		return appLogger
	}

	logger := AppLogger{}

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "ERROR"
	}
	logLevel = strings.ToUpper(logLevel)
	logsDir, ok := os.LookupEnv("LOGS_DIR")
	if !ok {
		panic("Incorrect/Missing value for LOGS_DIR")
	}

	zerolog.SetGlobalLevel(getLogLevel(logLevel))
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	if err := os.MkdirAll(logsDir, 0o744); err != nil {
		fmt.Println(fmt.Sprintf("error: %v", err))
	}
	filewriter := &lumberjack.Logger{
		Filename:   path.Join(logsDir, "apps-api.log"),
		MaxBackups: 4,
		MaxSize:    5,
		MaxAge:     10,
	}
	writers = append(writers, filewriter)
	// output := os.Stdout
	logWriters := io.MultiWriter(writers...)
	// domain -> product -> component -> module
	logger.logger = zerolog.New(logWriters).
		With().
		Str("Service", "apps-api").
		Timestamp().Logger()

	appLogger = &logger
	return appLogger
}

func getLogLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.ErrorLevel
	}
}

func (al *AppLogger) Trace(message string) {
	// add context as mandatory field
	al.logger.Trace().Msg(message)
}

func (al *AppLogger) Debug(message string) {
	// add context as mandatory field
	al.logger.Debug().Msg(message)
}

func (al *AppLogger) Info(message string) {
	// add context as mandatory field
	al.logger.Info().Msg(message)
}

func (al *AppLogger) Warn(message string) {
	// add context as mandatory field
	al.logger.Warn().Msg(message)
}

func (al *AppLogger) Error(err error, message string) {
	// add context as mandatory field
	al.logger.Error().Err(err).Msg(message)
}

func (al *AppLogger) Fatal(err error, message string) {
	// add context as mandatory field
	al.logger.Fatal().Err(err).Msg(message)
}

func (al *AppLogger) Panic(err error, message string) {
	// add context as mandatory field
	al.logger.Panic().Err(err).Msg(message)
}
