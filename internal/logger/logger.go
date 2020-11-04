package logger

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

const DefaultLogLevel = logrus.InfoLevel

// Entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {
	*logrus.Entry
}

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
func (l *Logger) Customize(cfg *Config) {
	l.config = *cfg
	l.SetLevel(l.config.Level)
}

// SetLevel injects log level to logger.
func (l *Logger) SetLevel(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = DefaultLogLevel
	}

	l.Logger.SetLevel(logrusLevel)
}

// SetService injects service name to logger.
func (l *Logger) SetService(service string) {
	l.service = service
}

// SetVersion injects service version to logger.
func (l *Logger) SetVersion(version string) {
	l.version = version
}

// WithAppInfo adds custom fields to the Entry.
func (l *Logger) WithAppInfo() *Entry {
	return &Entry{l.WithField("service", l.service).WithField("version", l.version)}
}
