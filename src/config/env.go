package config

import "os"

type Config struct {
	NetperfPort string
}

func Env() Config {
	netperfPort := os.Getenv("NETPERF_PORT")
	if netperfPort == "" {
		netperfPort = "12865"
	}

	return Config{
		NetperfPort: netperfPort,
	}
}

