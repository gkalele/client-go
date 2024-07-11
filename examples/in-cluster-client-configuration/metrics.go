package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"strings"
	"time"
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
	pingsSuccessful = generateNewCounter("pings_success")
	pingsFailed = generateNewCounter("pings_failed")
	dnsSuccessful = generateNewCounter("dns_success")
	dnsFailed = generateNewCounter("dns_failed")
	curlProbeSuccess = generateNewCounter("curl_probe_success")
	curlProbeFailed = generateNewCounter("curl_probe_failed")
}

func startPrometheusMetricsHandler(ctx context.Context) {
	http.Handle("/metrics", promhttp.Handler())
	go periodicMetricsPrinter(ctx)
	http.ListenAndServe(":2112", nil)
}

// periodicMetricsPrinter dumps all prometheus metrics every N minutes
func periodicMetricsPrinter(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			printMetrics()
		}
	}
}

// Copy of Gatherer.WriteTextToFile from prometheus - except we filter out and print only our metrics
func printMetrics() error {
	buf := &bytes.Buffer{}
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return err
	}
	for _, mf := range mfs {
		if strings.HasSuffix(*mf.Name, "failed") || strings.HasSuffix(*mf.Name, "success") {
			if _, err := expfmt.MetricFamilyToText(buf, mf); err != nil {
				return err
			}
		}
	}
	fmt.Println(buf.String())
	return nil
}
