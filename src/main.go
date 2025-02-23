package main

import (
	"time"

	"github.com/AposLaz/kube-netlag/config"
)

func main() {
    envVars := config.Env()

    if err := StartServer(envVars.NetperfPort); err != nil {
        panic(err)
    }

    nodes := GetTargetNodesIP()

    if len(nodes) == 0 {
        panic("No target nodes found.")
    }

    // declares a timer that will run every 5 seconds
    done := make(chan bool)

    for _, node := range nodes {
        go Monitoring(node,envVars.NetperfPort, done)
    }

    time.Sleep(120 * time.Second)
    
    // stop all go routines
    close(done)
    config.Logger("INFO", "All monitoring stopped.")
}