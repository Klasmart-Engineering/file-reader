package instrument

import (
	newrelic "github.com/newrelic/go-agent"
)

func NRES(txn newrelic.Transaction, url string) newrelic.ExternalSegment {
	return newrelic.ExternalSegment{
		StartTime: newrelic.StartSegmentNow(txn),
		URL:       url,
	}
}
