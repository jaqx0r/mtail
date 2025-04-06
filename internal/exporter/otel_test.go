// Copyright 2025 The mtail Authors. All Rights Reserved.
// This file is available under the Apache license.
// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"context"
	"testing"
	"time"

	"github.com/jaqx0r/mtail/internal/metrics"
	"github.com/jaqx0r/mtail/internal/metrics/datum"
	"github.com/jaqx0r/mtail/internal/testutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var otelTestCases = []struct {
	name      string
	progLabel bool
	metrics   []*metrics.Metric
	expected  []metricdata.ScopeMetrics
}{
	{
		name:    "empty",
		metrics: []*metrics.Metric{},
	},
	{
		name: "single",
		metrics: []*metrics.Metric{
			{
				Name:        "foo",
				Program:     "test",
				Kind:        metrics.Counter,
				Type:        metrics.Int,
				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
			},
		},
		expected: []metricdata.ScopeMetrics{{
			Scope: instrumentation.Scope{Name: "mtail_program"},
			Metrics: []metricdata.Metrics{
				{
					Name:        "foo",
					Description: "foo defined at ",
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Value: 1,
							},
						},
					},
				},
			},
		}},
	},
	{
		name:      "with prog label",
		progLabel: true,
		metrics: []*metrics.Metric{
			{
				Name:        "foo",
				Program:     "test",
				Kind:        metrics.Counter,
				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
			},
		},
		expected: []metricdata.ScopeMetrics{{
			Scope: instrumentation.Scope{Name: "mtail_program"},
			Metrics: []metricdata.Metrics{
				{
					Name:        "foo",
					Description: "foo defined at ",
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.String("prog", "test")),
								Value:      1,
							},
						},
					},
				},
			},
		}},
	},
		{
			name: "dimensioned",
			metrics: []*metrics.Metric{
				{
					Name:        "foo",
					Program:     "test",
					Kind:        metrics.Counter,
					Keys:        []string{"a", "b"},
					LabelValues: []*metrics.LabelValue{{Labels: []string{"1", "2"}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
				},
			},
		expected: []metricdata.ScopeMetrics{{
			Scope: instrumentation.Scope{Name: "mtail_program"},
			Metrics: []metricdata.Metrics{
				{
					Name:        "foo",
					Description: "foo defined at ",
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2")),
								Value:      1,
							},
						},
					},
				},
			},
		}},
		},
	//	{
	//		"gauge",
	//		false,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Gauge,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//			},
	//		},
	//		`# HELP foo defined at
	//
	// # TYPE foo gauge
	// foo{} 1
	// `,
	//
	//	},
	//	{
	//		"timer",
	//		false,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Timer,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//			},
	//		},
	//		`# HELP foo defined at
	//
	// # TYPE foo gauge
	// foo{} 1
	// `,
	//
	//	},
	//	{
	//		"text",
	//		false,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Text,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeString("hi", time.Unix(0, 0))}},
	//			},
	//		},
	//		"",
	//	},
	//	{
	//		"quotes",
	//		false,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Counter,
	//				Keys:        []string{"a"},
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{"str\"bang\"blah"}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//			},
	//		},
	//		`# HELP foo defined at
	//
	// # TYPE foo counter
	// foo{a="str\"bang\"blah"} 1
	// `,
	//
	//	},
	//	{
	//		"help",
	//		false,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Counter,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//				Source:      "location.mtail:37",
	//			},
	//		},
	//		`# HELP foo defined at location.mtail:37
	//
	// # TYPE foo counter
	// foo{} 1
	// `,
	//
	//	},
	//	{
	//		"2 help with label",
	//		true,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test2",
	//				Kind:        metrics.Counter,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//				Source:      "location.mtail:37",
	//			},
	//			{
	//				Name:        "foo",
	//				Program:     "test1",
	//				Kind:        metrics.Counter,
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{}, Value: datum.MakeInt(1, time.Unix(0, 0))}},
	//				Source:      "different.mtail:37",
	//			},
	//		},
	//		`# HELP foo defined at location.mtail:37
	//
	// # TYPE foo counter
	// foo{prog="test2"} 1
	// foo{prog="test1"} 1
	// `,
	//
	//	},
	//	{
	//		"histo",
	//		true,
	//		[]*metrics.Metric{
	//			{
	//				Name:        "foo",
	//				Program:     "test",
	//				Kind:        metrics.Histogram,
	//				Keys:        []string{"a"},
	//				LabelValues: []*metrics.LabelValue{{Labels: []string{"bar"}, Value: datum.MakeBuckets([]datum.Range{{0, 1}, {1, 2}}, time.Unix(0, 0))}},
	//				Source:      "location.mtail:37",
	//			},
	//		},
	//		`# HELP foo defined at location.mtail:37
	//
	// # TYPE foo histogram
	// foo_bucket{a="bar",prog="test",le="1"} 0
	// foo_bucket{a="bar",prog="test",le="2"} 0
	// foo_bucket{a="bar",prog="test",le="+Inf"} 0
	// foo_sum{a="bar",prog="test"} 0
	// foo_count{a="bar",prog="test"} 0
	// `,
	//
	//	},
	//	{
	//		"histo-count-eq-inf",
	//		true,
	//		[]*metrics.Metric{
	//			{
	//				Name:    "foo",
	//				Program: "test",
	//				Kind:    metrics.Histogram,
	//				Keys:    []string{"a"},
	//				LabelValues: []*metrics.LabelValue{
	//					{
	//						Labels: []string{"bar"},
	//						Value: &datum.Buckets{
	//							Buckets: []datum.BucketCount{
	//								{
	//									Range: datum.Range{Min: 0, Max: 1},
	//									Count: 1,
	//								},
	//								{
	//									Range: datum.Range{Min: 1, Max: 2},
	//									Count: 1,
	//								},
	//								{
	//									Range: datum.Range{Min: 2, Max: math.Inf(+1)},
	//									Count: 2,
	//								},
	//							},
	//							Count: 4,
	//							Sum:   5,
	//						},
	//					},
	//				},
	//				Source: "location.mtail:37",
	//			},
	//		},
	//		`# HELP foo defined at location.mtail:37
	//
	// # TYPE foo histogram
	// foo_bucket{a="bar",prog="test",le="1"} 1
	// foo_bucket{a="bar",prog="test",le="2"} 2
	// foo_bucket{a="bar",prog="test",le="+Inf"} 4
	// foo_sum{a="bar",prog="test"} 5
	// foo_count{a="bar",prog="test"} 4
	// `,
	//
	//	},
}

func TestOtelExport(t *testing.T) {
	for _, tc := range otelTestCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ms := metrics.NewStore()
			for _, metric := range tc.metrics {
				testutil.FatalIfErr(t, ms.Add(metric))
			}
			opts := []Option{
				Hostname("gunstar"),
			}
			if !tc.progLabel {
				opts = append(opts, OmitProgLabel())
			}
			e, err := New(ctx, ms, opts...)
			testutil.FatalIfErr(t, err)
			output, err := e.Produce(ctx)
			if err != nil {
				t.Error(err)
			}
			if len(output) != len(tc.expected) {
				t.Fatalf("output size doesn't match expected want %d got %d", len(tc.expected), len(output))
			}
			for i := range tc.expected {
				metricdatatest.AssertEqual(t, tc.expected[i], output[i], metricdatatest.IgnoreTimestamp())
			}
			e.Stop()
		})
	}
}
