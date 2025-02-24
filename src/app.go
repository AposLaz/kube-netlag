package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AposLaz/kube-netlag/config"
	"github.com/AposLaz/kube-netlag/k8s"
	"github.com/AposLaz/kube-netlag/netperf"
)

// GetTargetNodesIP returns a list of NodeInfo objects that represent the target nodes in the cluster
// for latency measurement. It will panic if it fails to create a Kubernetes client or fetch the
// cluster nodes.
func GetTargetNodesIP(currentNodeIp string) []k8s.NodeInfo {
	// get the client
	config, err := k8s.GetClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kubernetes client: %v", err))
	}

	// get the cluster nodes
	nodes, err := k8s.GetClusterNodes(config, currentNodeIp)
	if err != nil {
		panic(fmt.Sprintf("Failed to get cluster nodes: %v", err))
	}

	return nodes
}

var activeNodes sync.Map

// Monitoring starts a goroutine that continuously measures the latency to the target Node every 5 seconds.
// It will continue to run until the done channel is closed.
// If a node is already being monitored, the function will not start a new goroutine and will return.
// If an error occurs during the latency measurement, it will be logged, and the function will attempt to
// refresh the list of target nodes and restart the monitoring goroutine.
// The function will log the latency results for each node to the console.
func Monitoring(node k8s.NodeInfo, port string, currentNodeIp string, done <-chan bool) {
	// Check if the node is already being monitored
	if _, loaded := activeNodes.LoadOrStore(node.InternalIP, true); loaded {
		config.Logger("WARN", "Node %s is already being monitored. Skipping.", node.InternalIP)
		return
	}

	config.Logger("INFO", "Started monitoring Node: %s with IP: %s", node.Name, node.InternalIP)

	// Run the function every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		activeNodes.Delete(node.InternalIP) // Remove node from active tracking on exit
		config.Logger("INFO", "Stopped monitoring Node: %s with IP: %s", node.Name, node.InternalIP)
	}()

	for {
		select {
		case <-done:
			config.Logger("INFO", "Stopping latency monitoring for Node: %s with IP: %s", node.Name, node.InternalIP)
			return

		case <-ticker.C:
			config.Logger("INFO", "Monitoring Node: %s", node.Name)

			latency, err := netperf.ComputeLatency(node.InternalIP, port)
			if err != nil {
				config.Logger("ERROR", "Failed to compute latency for Node: %s with IP: %s\nError: %v", node.Name, node.InternalIP, err.Error())

				// Attempt to re-fetch node IPs
				config.Logger("INFO", "Attempting to refresh node IPs...")
				newNodes := GetTargetNodesIP(currentNodeIp)
				if len(newNodes) == 0 {
					config.Logger("ERROR", "Failed to refresh node IPs. Retrying next cycle.")
					continue
				}

				config.Logger("INFO", "Successfully refreshed node IPs. Restarting monitoring.")

				// Start new monitoring goroutines for refreshed nodes
				for _, newNode := range newNodes {
					go Monitoring(newNode, port, currentNodeIp, done)
				}

				// Stop the current goroutine to avoid duplicates
				return
			}

			config.Logger("INFO", "Latency Results | from_node=%s current_ip=%s to_node=%s target_ip=%s min_latency_ms=%.2f max_latency_ms=%.2f mean_latency_ms=%.2f",
				node.CurrentNodeName, currentNodeIp, node.Name, node.InternalIP, latency[0], latency[1], latency[2])
		}
	}
}

// StartServer initializes and starts the netperf server on the specified port.
// It delegates the actual startup process to the netperf.StartServer function,
// which handles retries and logs the outcome. Returns an error if the server
// fails to start after the maximum number of retries.
func StartServer(port string) error {
	return netperf.StartServer(port)
}
