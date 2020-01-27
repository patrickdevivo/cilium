// Copyright 2019-2020 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	registry    *prometheus.Registry
	ApiDuration *prometheus.HistogramVec
	RateLimit   *prometheus.HistogramVec
}

// NewLegacyPrometheusMetrics returns a new metrics tracking implementation to
// cover external API usage. This is the legacy version of the function to
// allow providing backward compatibility to existing metric names. For new
// code, use NewPrometheusMetrics.
func NewLegacyPrometheusMetrics(namespace, subsystem, durationName, rateLimitName string, registry *prometheus.Registry) *prometheusMetrics {
	m := &prometheusMetrics{
		registry: registry,
	}

	m.ApiDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      durationName,
		Help:      "Duration of interactions with API",
	}, []string{"operation", "responseCode"})

	m.RateLimit = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      rateLimitName,
		Help:      "Duration of client-side rate limiter blocking",
	}, []string{"operation"})

	registry.MustRegister(m.ApiDuration)
	registry.MustRegister(m.RateLimit)

	return m
}

// NewPrometheusMetrics returns a new metrics tracking implementation to cover
// external API usage.
func NewPrometheusMetrics(namespace, subsystem string, registry *prometheus.Registry) *prometheusMetrics {
	return NewLegacyPrometheusMetrics(namespace, subsystem, "api_duration_seconds", "rate_limit_duration_seconds", registry)
}

func (p *prometheusMetrics) ObserveAPICall(operation, status string, duration float64) {
	p.ApiDuration.WithLabelValues(operation, status).Observe(duration)
}

func (p *prometheusMetrics) ObserveRateLimit(operation string, delay time.Duration) {
	p.RateLimit.WithLabelValues(operation).Observe(delay.Seconds())
}

// NoOpMetrics is a no-op implementation
type NoOpMetrics struct{}

func (m *NoOpMetrics) ObserveAPICall(call, status string, duration float64)      {}
func (m *NoOpMetrics) ObserveRateLimit(operation string, duration time.Duration) {}
