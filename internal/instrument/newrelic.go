package instrument

import (
	"strconv"

	"github.com/KL-Engineering/file-reader/internal/log"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrzap"
)

type NewRelic struct {
	serviceName string
	cfg         newrelic.Config
	app         newrelic.Application
}

func GetNewRelic(serviceName string, logger *log.ZapLogger) (*NewRelic, error) {
	cfg := newrelic.NewConfig(serviceName, MustGetEnv("NEW_RELIC_LICENSE_KEY"))
	isDistributedTracerEnabled := false
	isSpanEventsEnabled := false
	isErrorCollectorEnabled := false

	isDistributedTracerEnabled, _ = strconv.ParseBool(MustGetEnv("DISTRIBUTED_TRACER_ENABLED"))
	isSpanEventsEnabled, _ = strconv.ParseBool(MustGetEnv("SPAN_EVENT_ENABLED"))
	isErrorCollectorEnabled, _ = strconv.ParseBool(MustGetEnv("ERROR_COLLECTOR_ENABLED"))

	cfg.DistributedTracer.Enabled = isDistributedTracerEnabled
	cfg.SpanEvents.Enabled = isSpanEventsEnabled
	cfg.ErrorCollector.Enabled = isErrorCollectorEnabled
	newRelicEnabled, _ := strconv.ParseBool(MustGetEnv("NEW_RELIC_ENABLED"))
	cfg.Enabled = newRelicEnabled
	cfg.Logger = nrzap.Transform(logger.Named("newrelic"))

	app, err := newrelic.NewApplication(cfg)

	if err != nil {
		panic("Failed to setup NewRelic: " + err.Error())
	}
	return &NewRelic{
		serviceName: serviceName,
		cfg:         cfg,
		app:         app,
	}, nil
}
