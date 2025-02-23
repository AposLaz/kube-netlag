package app

import (
	"time"

	"github.com/AposLaz/kube-netlag/config"
	"github.com/AposLaz/kube-netlag/netperf"
)

func Monitoring(ip string, done <- chan bool){
	// run the function every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	// this function will stop the ticker when the program exits (avoid memory leaks)
	defer ticker.Stop()

	for {
		select {
		case <- done:
			config.Logger("INFO", "Stopping latency monitoring for Node: %s", ip)
		case <- ticker.C:
			config.Logger("INFO", "Monitoring Node: %s", ip)
                
            latency, err := netperf.ComputeLatency(ip)
            if err != nil {
                config.Logger("ERROR", "Failed to compute latency for Node: %s\nError: %v", ip, err.Error())
                continue
            }

			config.Logger("INFO", "Latency Results | Node: %s | Min: %.2f ms | Max: %.2f ms | Mean: %.2f ms",
				ip, latency[0], latency[1], latency[2])
		}
	}
}