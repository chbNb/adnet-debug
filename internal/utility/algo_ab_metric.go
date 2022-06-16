package utility

import (
	"fmt"

	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type AlgoABTestMetric struct{}

func NewAlgoABTestMetric() *AlgoABTestMetric {
	return new(AlgoABTestMetric)
}

func (a *AlgoABTestMetric) ErrorRecord(where int, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
		if len(msg) > 24 {
			msg = msg[:24]
		}
	}
	metrics.IncCounterWithLabelValues(30, fmt.Sprintf("%d_%s", where, msg))
}
