package logger

// Config contains Logger settings.
type Config struct {
	Level string `yaml:"level"`
}

// Option allows to inject options to Logger.
type Option func(l *Logger)

// SetLevel injects log level to logger.
func SetLevel(level string) Option {
	return func(l *Logger) {
		l.SetLevel(level)
	}
}

// SetLSetServiceevel injects service name to logger.
func SetService(service string) Option {
	return func(l *Logger) {
		l.SetService(service)
	}
}

// SetVersion injects service version to logger.
func SetVersion(version string) Option {
	return func(l *Logger) {
		l.SetVersion(version)
	}
}
