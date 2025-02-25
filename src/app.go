package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AposLaz/kube-netlag/config"
	"github.com/AposLaz/kube-netlag/k8s"
	"github.com/AposLaz/kube-netlag/netperf"
	"github.com/AposLaz/kube-netlag/promMetrics"
)

type CurrentNodeInfo struct {
	Name       string
	InternalIP string
}

// GetTargetNodesIP returns a list of NodeInfo objects that represent the target nodes in the cluster
// for latency measurement. It will panic if it fails to create a Kubernetes client or fetch the
// cluster nodes.
func GetTargetNodesIP(currentNodeIp string) (string, []k8s.NodeInfo) {
	// get the client
	config, err := k8s.GetClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kubernetes client: %v", err))
	}

	// get the cluster nodes
	currentNode, nodes, err := k8s.GetClusterNodes(config, currentNodeIp)
	if err != nil {
		panic(fmt.Sprintf("Failed to get cluster nodes: %v", err))
	}

	return currentNode, nodes
}

var activeNodes sync.Map

// MonitoringLatency initiates a latency monitoring process for a given node.
// It periodically computes the latency from the current node to the target node
// using the netperf tool and updates Prometheus metrics with the results.
// The monitoring runs in a separate goroutine and continues until the node is
// either removed from the monitoring list or an error occurs.
//
// Parameters:
//
//	node: The target node to monitor, including its name and internal IP address.
//	port: The port number on which the netperf server is running on the target node.
//	currentNode: Information about the current node (name and internal IP).
//	failureChan: A channel to report monitoring failures, where the nodes IP is sent in case of failure.
//
// The function ensures that a node is not monitored multiple times concurrently by
// checking and updating the activeNodes map. It logs the start and stop of monitoring,
// as well as any errors encountered during latency computation. The monitoring is
// interrupted if an error occurs, with the node's IP sent through the failureChan.
func MonitoringLatency(node k8s.NodeInfo, port string, currentNode CurrentNodeInfo, failureChan chan<- string) {
	// Check if the node is already being monitored
	if _, loaded := activeNodes.LoadOrStore(node.InternalIP, true); loaded {
		config.Logger("WARN", "Node %s is already being monitored. Skipping.", node.InternalIP)
		return
	}

	config.Logger("INFO", "Started monitoring Node: %s with IP: %s", node.Name, node.InternalIP)

	defer func() {
		activeNodes.Delete(node.InternalIP)
		config.Logger("INFO", "Stopped monitoring Node: %s with IP: %s", node.Name, node.InternalIP)
	}()

	for {
		config.Logger("INFO", "Monitoring Node: %s", node.Name)

		latency, err := netperf.ComputeLatency(node.InternalIP, port)
		if err != nil {
			config.Logger("ERROR", "Failed to compute latency for Node: %s with IP: %s\nError: %v", node.Name, node.InternalIP, err.Error())
			failureChan <- node.InternalIP // Report failure to main
			return
		}

		config.Logger("INFO", "Latency Results | from_node=%s current_ip=%s to_node=%s target_ip=%s min_latency_ms=%.2f max_latency_ms=%.2f mean_latency_ms=%.2f",
			currentNode.Name, currentNode.InternalIP, node.Name, node.InternalIP, latency[0], latency[1], latency[2])

		metrics := promMetrics.LatencyMeasurement{FromNodeName: currentNode.Name, FromIpAddress: currentNode.InternalIP, ToNodeName: node.Name, ToIpAddress: node.InternalIP, MinLatency: latency[0], MaxLatency: latency[1], AvgLatency: latency[2]}
		promMetrics.UpdateMetrics(metrics)

		time.Sleep(10 * time.Second)
	}
}

// NetperfServer launches the netperf server on the specified port. It attempts to start the server
// up to a maximum number of retries if initial attempts fail. The function logs the success or
// failure of starting the server and returns an error if all attempts are unsuccessful.
func StartNetperfServer(port string) error {
	return netperf.StartServer(port)
}

// InitializeMonitoring starts the monitoring process for the given environment variables.
// It fetches the target nodes in the cluster, starts a goroutine to monitor each target node,
// and then enters a loop to refresh the target nodes and handle any failed nodes.
// The loop exits when the process receives an interrupt or termination signal.
func InitializeMonitoring(envVars config.EnvVars) {
	currentNode, nodes := GetTargetNodesIP(envVars.CurrentNodeIp)
	if len(nodes) == 0 || currentNode == "" {
		panic("No target nodes found.")
	}

	failureChan := make(chan string)
	currentNodeInfo := CurrentNodeInfo{Name: currentNode, InternalIP: envVars.CurrentNodeIp}

	for _, node := range nodes {
		if node.InternalIP != envVars.CurrentNodeIp {
			go MonitoringLatency(node, envVars.NetperfPort, currentNodeInfo, failureChan)
		}
	}

	refreshTicker := time.NewTicker(1 * time.Minute)
	defer refreshTicker.Stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case <-refreshTicker.C:
			handleNodeRefresh(envVars, failureChan)
		case failedIP := <-failureChan:
			handleNodeFailure(envVars, failedIP, failureChan)
		case <-signals:
			run = false
			config.Logger("INFO", "Shutting down monitoring...")
		}
	}

	time.Sleep(5 * time.Second)
}

// handleNodeRefresh updates the monitoring state of nodes in the cluster.
// It fetches the current list of nodes and compares it with the actively monitored nodes.
// If a node is new and not the current node, it starts monitoring latency for that node.
// Nodes that are no longer part of the cluster are removed from the active monitoring map.
func handleNodeRefresh(envVars config.EnvVars, failureChan chan<- string) {
	currentNode, newNodes := GetTargetNodesIP(envVars.CurrentNodeIp)
	existingNodes := make(map[string]bool)

	currentNodeInfo := CurrentNodeInfo{Name: currentNode, InternalIP: envVars.CurrentNodeIp}

	for _, node := range newNodes {
		existingNodes[node.InternalIP] = true
		if _, loaded := activeNodes.Load(node.InternalIP); !loaded && node.InternalIP != envVars.CurrentNodeIp {
			go MonitoringLatency(node, envVars.NetperfPort, currentNodeInfo, failureChan)
		}
	}

	activeNodes.Range(func(key, value interface{}) bool {
		ip := key.(string)
		if !existingNodes[ip] {
			activeNodes.Delete(ip)
			config.Logger("INFO", "Node %s removed from monitoring due to cluster update.", ip)
		}
		return true
	})
}

// handleNodeFailure restarts the monitoring process for a node that has previously failed.
// It logs the restart action and retrieves the current list of nodes. If the failed node is
// still present in the cluster, it initiates the MonitoringLatency function for that node
// using the provided netperf port and current node information.
func handleNodeFailure(envVars config.EnvVars, failedIP string, failureChan chan<- string) {
	config.Logger("INFO", "Restarting monitoring for Node with IP: %s", failedIP)
	currentNode, newNodes := GetTargetNodesIP(envVars.CurrentNodeIp)

	currentNodeInfo := CurrentNodeInfo{Name: currentNode, InternalIP: envVars.CurrentNodeIp}

	for _, node := range newNodes {
		if node.InternalIP == failedIP {
			go MonitoringLatency(node, envVars.NetperfPort, currentNodeInfo, failureChan)
		}
	}
}
