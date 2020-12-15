package logger

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// fields
	ErrorField = "error"
	TitleField = "title"

	// outputs
	Stdout = "stdout"
	Stderr = "stderr"
	Vacuum = "vacuum"

	// formatters
	JSONFormatter = "json"
	TextFormatter = "text"
)

var formatters = map[string]logrus.Formatter{
	JSONFormatter: &logrus.JSONFormatter{},
	TextFormatter: &logrus.TextFormatter{},
}

var levelMap = map[string]logrus.Level{
	"panic":   logrus.PanicLevel,
	"fatal":   logrus.FatalLevel,
	"error":   logrus.ErrorLevel,
	"warning": logrus.WarnLevel,
	"info":    logrus.InfoLevel,
	"debug":   logrus.DebugLevel,
	"trace":   logrus.TraceLevel,
}

// Config is a configuration for logger
type Config struct {
	Output    string `yaml:"output" json:"output" toml:"output"`          // enum (stdout|stderr|vacuum|path/to/file)
	Formatter string `yaml:"formatter" json:"formatter" toml:"formatter"` // enum (json|text)
	Level     string `yaml:"level" json:"level" toml:"level"`             // enum (panic|fatal|error|warning|info|debug|trace)
}

// CreateLogger returns logger acording config
func (config Config) CreateLogger() *logrus.Logger {
	logger := logrus.New()

	switch config.Output {
	case Stdout, "":
		logger.SetOutput(os.Stdout)
	case Stderr:
		logger.SetOutput(os.Stderr)
	case Vacuum:
		logger.SetOutput(ioutil.Discard)
	default:
		f, err := os.OpenFile(config.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			logger.SetOutput(os.Stdout)
			logger.WithError(err).Debug("falling to stdout")
		} else {
			logger.SetOutput(f)
		}
	}

	formatter, ok := formatters[config.Formatter]
	if !ok {
		formatter = &logrus.TextFormatter{}
	}
	logger.SetFormatter(formatter)

	lvl, ok := levelMap[config.Level]
	if !ok {
		lvl = logrus.DebugLevel
	}
	logger.SetLevel(lvl)

	return logger
}
