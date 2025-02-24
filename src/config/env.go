package config

import "os"

type Config struct {
	NetperfPort   string
	CurrentNodeIp string
	MetricsPort   string
}

// Env returns a Config object with environment variable values. If a variable is
// unset, it will use the following default values:
// - NETPERF_PORT: 12865
// - METRICS_PORT: 9090
// - HOST_IP: "" (must be set)
func Env() Config {
	netperfPort := os.Getenv("NETPERF_PORT")
	if netperfPort == "" {
		netperfPort = "12865"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	return Config{
		NetperfPort:   netperfPort,
		CurrentNodeIp: os.Getenv("HOST_IP"),
		MetricsPort:   metricsPort,
	}
}
