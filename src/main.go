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
	go promMetrics.StartServer(envVars.MetricsPort)

	if err := StartNetperfServer(envVars.NetperfPort); err != nil {
		panic(err)
	}

	InitializeMonitoring(envVars)
}
