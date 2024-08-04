package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	InfoLogger  zerolog.Logger
	ErrorLogger zerolog.Logger
)

func init() {
	stdoutWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	stdoutLogger := zerolog.New(stdoutWriter).With().Timestamp().Logger()
	stderrWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	stderrLogger := zerolog.New(stderrWriter).With().Timestamp().Logger()
	InfoLogger = stdoutLogger
	ErrorLogger = stderrLogger
}
