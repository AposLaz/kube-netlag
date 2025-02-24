package config

import "os"

type Config struct {
	NetperfPort   string
	CurrentNodeIp string
	MetricsPort   string
}

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
