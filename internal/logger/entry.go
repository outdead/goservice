package logger

import "github.com/sirupsen/logrus"

// Entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {
	*logrus.Entry
	logger *Logger
}

// Logger returns *Logger instance.
func (e *Entry) Logger() *Logger {
	return e.logger
}
