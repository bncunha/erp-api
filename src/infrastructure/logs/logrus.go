package logs

import "github.com/sirupsen/logrus"

type logrusLogger struct {
	entry *logrus.Entry
}

func NewLogrus() Logs {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	entry := logrus.NewEntry(logger)
	return &logrusLogger{
		entry: entry,
	}
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logrusLogger) AddHook(h Hook) {
	l.entry.Logger.AddHook(&logrusHookAdapter{h: h})
}

func (l *logrusLogger) With(fields map[string]any) Logs {
	return &logrusLogger{
		entry: l.entry.WithFields(fields),
	}
}

type logrusHookAdapter struct {
	h Hook
}

func (a *logrusHookAdapter) Levels() []logrus.Level {
	levels := []logrus.Level{}
	for _, lvl := range a.h.Levels() {
		switch lvl {
		case "info":
			levels = append(levels, logrus.InfoLevel)
		case "error":
			levels = append(levels, logrus.ErrorLevel)
		default:
			levels = append(levels, logrus.InfoLevel)
		}
	}
	return levels
}

func (a *logrusHookAdapter) Fire(entry *logrus.Entry) error {
	// converte para map[string]any e chama o hook genérico
	data := map[string]any{}
	for k, v := range entry.Data {
		data[k] = v
	}
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()
	return a.h.Fire(data)
}
