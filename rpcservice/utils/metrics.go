// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var mt metric.Meter

func init() {
	meterProvider := otel.GetMeterProvider()
	mt = meterProvider.Meter("service")
}

func NewFloat64Histogram(metricName string, tagNames ...string) metric.Float64Histogram {
	ht, _ := mt.Float64Histogram(metricName)
	return ht
}

func NewInt64Counter(metricName string, tagNames ...string) metric.Int64Counter {
	ct, _ := mt.Int64Counter(metricName)
	return ct
}

func TimeSinceMicroseconds(startTime time.Time) int {
	return int(time.Since(startTime).Microseconds())
}

func EmitThroughputAndLatency(ht metric.Float64Histogram, ct metric.Int64Counter, startTime time.Time, attr ...attribute.KeyValue) {
	EmitLatency(ht, startTime)
	ct.Add(context.Background(), 1, metric.WithAttributes(attr...))
}

func EmitLatency(ht metric.Float64Histogram, startTime time.Time, attr ...attribute.KeyValue) {
	latency := TimeSinceMicroseconds(startTime)
	elapsedTime := float64(latency) / float64(time.Millisecond)
	ht.Record(context.Background(), elapsedTime, metric.WithAttributes(attr...))
}
