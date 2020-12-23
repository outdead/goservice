package logutils

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

const DefaultLogLevel = logrus.InfoLevel

// Logger contains wrapped Logrus logger.
type Logger struct {
	*logrus.Logger
	config  Config
	service string
	version string
}

// New creates new logger.
func New(options ...Option) *Logger {
	logger := Logger{
		Logger: logrus.New(),
	}
	logger.Formatter = new(logrus.JSONFormatter)

	for _, option := range options {
		option(&logger)
	}

	return &logger
}

// NewDiscardLogger creates discard Logger on which all Write calls succeed
// without doing anything.
func NewDiscardLogger() *Logger {
	logger := New()
	logger.Out = ioutil.Discard

	return logger
}

// Customize gets values from config and customizes Logger.
func (log *Logger) Customize(cfg *Config) {
	log.config = *cfg
	log.SetLevel(log.config.Level)
}

// SetLevel injects log level to logger.
func (log *Logger) SetLevel(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = DefaultLogLevel
	}

	log.Logger.SetLevel(logrusLevel)
}

// SetService injects service name to logger.
func (log *Logger) SetService(service string) {
	log.service = service
}

// SetVersion injects service version to logger.
func (log *Logger) SetVersion(version string) {
	log.version = version
}

// NewEntry returns new empty entry.
func (log *Logger) NewEntry() *Entry {
	entry := logrus.NewEntry(log.Logger)

	if log.service != "" {
		entry = entry.WithField("service", log.service)
	}

	if log.version != "" {
		entry = entry.WithField("version", log.version)
	}

	return &Entry{Entry: entry, logger: log}
}

// Version returns injected service version.
func (log *Logger) Version() string {
	return log.version
}
