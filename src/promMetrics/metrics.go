package promMetrics

import (
	"net/http"

	"github.com/AposLaz/kube-netlag/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LatencyMeasurement struct {
	FromNodeName  string
	FromIpAddress string
	ToNodeName    string
	ToIpAddress   string
	MinLatency    float64
	MaxLatency    float64
	AvgLatency    float64
}

var (
	// Define Prometheus Gauges for latency metrics
	minLatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_min_latency_ms",
			Help: "Minimum latency in milliseconds between nodes.",
		},
		[]string{"from_node", "to_node", "from_ip", "to_ip"},
	)

	maxLatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_max_latency_ms",
			Help: "Maximum latency in milliseconds between nodes.",
		},
		[]string{"from_node", "to_node", "from_ip", "to_ip"},
	)

	avgLatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_avg_latency_ms",
			Help: "Average latency in milliseconds between nodes.",
		},
		[]string{"from_node", "to_node", "from_ip", "to_ip"},
	)
)

// Init registers the Prometheus gauges for latency metrics with the default
// registry. It should be called once at application startup to enable
// Prometheus metrics collection.
func Init() {
	prometheus.MustRegister(minLatencyGauge)
	prometheus.MustRegister(maxLatencyGauge)
	prometheus.MustRegister(avgLatencyGauge)
}

// UpdateMetrics updates the Prometheus gauges with the given latency metrics.
func UpdateMetrics(metrics LatencyMeasurement) {
	labels := prometheus.Labels{
		"from_node": metrics.FromNodeName,
		"to_node":   metrics.ToNodeName,
		"from_ip":   metrics.FromIpAddress,
		"to_ip":     metrics.ToIpAddress,
	}

	minLatencyGauge.With(labels).Set(metrics.MinLatency)
	maxLatencyGauge.With(labels).Set(metrics.MaxLatency)
	avgLatencyGauge.With(labels).Set(metrics.AvgLatency)
}

// StartServer initializes an HTTP server on the specified port to expose Prometheus metrics.
// It registers the "/metrics" endpoint and starts listening for incoming requests.
// If the server fails to start, it logs an error message and panics.
func StartServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		config.Logger("ERROR", "Failed to start prometheus server: %v", err)
		panic(err)
	}
}
