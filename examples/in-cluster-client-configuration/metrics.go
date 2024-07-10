package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	pingsSuccessful  prometheus.Counter
	pingsFailed      prometheus.Counter
	dnsSuccessful    prometheus.Counter
	dnsFailed        prometheus.Counter
	curlProbeSuccess prometheus.Counter
	curlProbeFailed  prometheus.Counter
)

func generateNewCounter(name string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{Name: name, Help: name})
}

func init() {
	pingsSuccessful = generateNewCounter("pings.success")
	pingsFailed = generateNewCounter("pings.failed")
	dnsSuccessful = generateNewCounter("dns.success")
	dnsFailed = generateNewCounter("dns.failed")
	curlProbeSuccess = generateNewCounter("curl.probe.success")
	curlProbeFailed = generateNewCounter("curl.probe.failed")
}

func startPrometheusMetricsHandler() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
