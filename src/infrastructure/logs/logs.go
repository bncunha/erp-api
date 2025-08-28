package logs

type Hook interface {
	Levels() []string // use string para não acoplar ao logrus.Level
	Fire(entry map[string]any) error
}

type Logs interface {
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	AddHook(h Hook)
	With(fields map[string]any) Logs
}

var Logger Logs

func NewLogs() {
	if Logger != nil {
		return
	}
	loggerApp := NewLogrus()
	Logger = loggerApp
}
