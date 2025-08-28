package logs

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type nrHook struct{ app *newrelic.Application }

func NewRelicHook(app *newrelic.Application) Hook { return &nrHook{app} }

func (h *nrHook) Levels() []string {
	allLevels := logrus.AllLevels
	levels := []string{}
	for _, lvl := range allLevels {
		levels = append(levels, lvl.String())
	}
	return levels
}

func (h *nrHook) Fire(entry map[string]any) error {
	// opcional: transformar em evento/atributos
	if h.app == nil {
		return nil
	}
	h.app.RecordCustomEvent("LogEvent", toAttrs(entry))
	return nil
}

func toAttrs(entry map[string]any) map[string]interface{} {
	e := entry["logrus"].(*logrus.Entry)
	m := map[string]interface{}{
		"level":   e.Level.String(),
		"message": e.Message,
	}
	for k, v := range e.Data {
		m[k] = v
	}
	return m
}
