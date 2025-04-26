// Copyright 2025 The mtail Authors.  All Rights Reserved
// This file is available under the Apache license.

package exporter

import (
	"context"
	"fmt"
	"math"
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
		case metrics.Gauge:
			switch m.Type {
			case metrics.Int:
				newMetric.Data = otelIntGauge(m, e.omitProgLabel)
			case metrics.Float:
				newMetric.Data = otelFloatGauge(m, e.omitProgLabel)
			default:
				return nil
			}
		case metrics.Timer:
			switch m.Type {
			case metrics.Float:
				newMetric.Data = otelFloatGauge(m, e.omitProgLabel)
			default:
				return nil
			}
		case metrics.Histogram:
			switch m.Type {
			case metrics.Buckets:
				newMetric.Data = otelHisto(m, e.omitProgLabel)
			default:
				return nil
			}
		default:
			return nil
		}
		otelMetrics = append(otelMetrics, newMetric)
		return nil
	})

	if len(otelMetrics) == 0 {
		return nil, nil
	}

	return []metricdata.ScopeMetrics{{
		Scope:   instrumentation.Scope{Name: "mtail_program"},
		Metrics: otelMetrics,
	}}, nil
}

func otelIntCounter(m *metrics.Metric, omitProgLabel bool) metricdata.Sum[int64] {
	counter := metricdata.Sum[int64]{
		DataPoints:  make([]metricdata.DataPoint[int64], 0, len(m.LabelValues)),
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

func otelIntGauge(m *metrics.Metric, omitProgLabel bool) metricdata.Gauge[int64] {
	gauge := metricdata.Gauge[int64]{
		DataPoints: make([]metricdata.DataPoint[int64], 0, len(m.LabelValues)),
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
		gauge.DataPoints = append(gauge.DataPoints, dp)
	}
	return gauge
}

func otelFloatGauge(m *metrics.Metric, omitProgLabel bool) metricdata.Gauge[float64] {
	gauge := metricdata.Gauge[float64]{
		DataPoints: make([]metricdata.DataPoint[float64], 0, len(m.LabelValues)),
	}
	lsc := make(chan *metrics.LabelSet)
	go m.EmitLabelSets(lsc)
	for ls := range lsc {
		dp := metricdata.DataPoint[float64]{
			Attributes: otelLabels(ls.Labels, omitProgLabel, m.Program),
			StartTime:  processStartTime,
			Time:       ls.Datum.TimeUTC(),
			Value:      datum.GetFloat(ls.Datum),
		}
		gauge.DataPoints = append(gauge.DataPoints, dp)
	}
	return gauge
}

func otelHisto(m *metrics.Metric, omitProgLabel bool) metricdata.Histogram[float64] {
	histo := metricdata.Histogram[float64]{
		DataPoints:  make([]metricdata.HistogramDataPoint[float64], 0, len(m.LabelValues)),
		Temporality: metricdata.CumulativeTemporality,
	}
	lsc := make(chan *metrics.LabelSet)
	go m.EmitLabelSets(lsc)
	for ls := range lsc {
		bounds, counts := otelConvertBuckets(datum.GetBuckets(ls.Datum))
		dp := metricdata.HistogramDataPoint[float64]{
			Attributes:   otelLabels(ls.Labels, omitProgLabel, m.Program),
			StartTime:    processStartTime,
			Time:         ls.Datum.TimeUTC(),
			Count:        datum.GetBucketsCount(ls.Datum),
			Sum:          datum.GetBucketsSum(ls.Datum),
			Bounds:       bounds,
			BucketCounts: counts,
		}
		histo.DataPoints = append(histo.DataPoints, dp)
	}
	return histo
}

func otelConvertBuckets(d *datum.Buckets) (bounds []float64, counts []uint64) {
	if len(d.Buckets) == 0 {
		// Should never happen?
		return nil, nil
	}
	// The last bucket may be the +Inf bucket, which is implied in OTel, but explicit in mtail.
	if math.IsInf(d.Buckets[len(d.Buckets)-1].Range.Max, +1) {
		bounds = make([]float64, len(d.Buckets)-1)
	} else {
		bounds = make([]float64, len(d.Buckets))
	}
	counts = make([]uint64, len(d.Buckets))
	for i, bucket := range d.Buckets {
		if bound := bucket.Range.Max; !math.IsInf(bound, +1) {
			bounds[i] = bound
		}
		counts[i] = bucket.Count
	}
	return
}

func otelLabels(labels map[string]string, omitProgLabel bool, programName string) attribute.Set {
	l := len(labels)
	if !omitProgLabel {
		l++
	}
	kvs := make([]attribute.KeyValue, l)
	i := 0
	for k, v := range labels {
		kvs[i] = attribute.String(k, v)
		i++
	}
	if !omitProgLabel {
		kvs[i] = attribute.String("prog", programName)
	}
	return attribute.NewSet(kvs...)
}
