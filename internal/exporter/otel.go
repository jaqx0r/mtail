// Copyright 2025 The mtail Authors.  All Rights Reserved
// This file is available under the Apache license.

package exporter

import (
	"context"
	"fmt"
	"time"

	"github.com/jaqx0r/mtail/internal/metrics"
	"github.com/jaqx0r/mtail/internal/metrics/datum"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var processStartTime = time.Now()

// Produce implements the opentelemetry Producer.Produce method.
func (e *Exporter) Produce(context.Context) ([]metricdata.ScopeMetrics, error) {
	otelMetrics := make([]metricdata.Metrics, 0)

	e.store.Range(func(m *metrics.Metric) error {
		m.RLock()
		defer m.RUnlock()

		newMetric := metricdata.Metrics{
			Name:        m.Name,
			Description: fmt.Sprintf("%s defined at %s", m.Name, m.Source),
		}
		switch m.Kind {
		case metrics.Counter:
			switch m.Type {
			case metrics.Int:
				newMetric.Data = otelIntCounter(m, e.omitProgLabel)
			default:
				return nil
			}
		default:
			return nil
		}
		otelMetrics = append(otelMetrics, newMetric)
		return nil
	})

	return []metricdata.ScopeMetrics{{
		Scope:   instrumentation.Scope{Name: "mtail_program"},
		Metrics: otelMetrics,
	}}, nil
}

func otelIntCounter(m *metrics.Metric, omitProgLabel bool) metricdata.Sum[int64] {
	counter := metricdata.Sum[int64]{
		DataPoints:  make([]metricdata.DataPoint[int64], 0),
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
	}
	lsc := make(chan *metrics.LabelSet)
	go m.EmitLabelSets(lsc)
	for ls := range lsc {
		dp := metricdata.DataPoint[int64]{
			Attributes: otelLabels(ls.Labels, omitProgLabel, m.Program),
			StartTime:  processStartTime,
			Time:       ls.Datum.TimeUTC(),
			Value:      datum.GetInt(ls.Datum),
		}
		counter.DataPoints = append(counter.DataPoints, dp)
	}
	return counter

}

func otelLabels(labels map[string]string, omitProgLabel bool, programName string) attribute.Set {
	kvs := make([]attribute.KeyValue, len(labels))
	i := 0
	for k, v := range labels {
		kvs[i] = attribute.String(k, v)
		i++
	}
	if omitProgLabel {
		kvs[i] = attribute.String("prog", programName)
	}
	return attribute.NewSet(kvs...)
}
