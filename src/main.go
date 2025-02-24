package main

import (
	"github.com/AposLaz/kube-netlag/config"
	"github.com/AposLaz/kube-netlag/promMetrics"
)

func main() {
	envVars := config.Env()

	// intialize prometheus metrics
	promMetrics.Init()
	// Initialize prometheus server
	promMetrics.StartServer(envVars.MetricsPort)

	if err := StartServer(envVars.NetperfPort); err != nil {
		panic(err)
	}

	nodes := GetTargetNodesIP(envVars.CurrentNodeIp)

	if len(nodes) == 0 {
		panic("No target nodes found.")
	}

	// declares a timer that will run every 5 seconds
	done := make(chan bool)

	for _, node := range nodes {
		go Monitoring(node, envVars.NetperfPort, envVars.CurrentNodeIp, done)
	}

	select {}
}
