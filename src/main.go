package main

import (
	"time"

	"github.com/AposLaz/kube-netlag/app"
	"github.com/AposLaz/kube-netlag/config"
)

func main() {
    envVars := config.Env()

    if err := app.StartServer(envVars.NetperfPort); err != nil {
        panic(err)
        return
    }
    // TODO fetch the ips address from the Nodes
    ips := []string{"0.0.0.0","127.0.0.1", "244.178.44.111"}

    // declares a timer that will run every 5 seconds
    done := make(chan bool)

    for _, ip := range ips {
        go app.Monitoring(ip,envVars.NetperfPort, done)
    }

    time.Sleep(120 * time.Second)
    
    // stop all go routines
    close(done)
    config.Logger("INFO", "All monitoring stopped.")
}