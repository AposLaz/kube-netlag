package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeInfo struct {
	Name       string
	ExternalIP string
}

func GetClusterNodes(clientset *kubernetes.Clientset, currentNodeIP string) ([]NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list nodes: %w", err)
	}

	var nodesInfo []NodeInfo
	for _, node := range nodes.Items {
		var externalIP string

		for _, addr := range node.Status.Addresses {
			if addr.Type == "ExternalIP" {
				externalIP = addr.Address
			}
		}

		if externalIP == currentNodeIP {
			continue
		}

		nodesInfo = append(nodesInfo, NodeInfo{
			Name:       node.Name,
			ExternalIP: externalIP,
		})
	}

	return nodesInfo, nil
}

// GetCurrentNodeIP returns the IP of the current node, as determined by running `nslookup kubernetes.default.svc.cluster.local`.
// If the command fails, or if the current node's IP is not found, an error is returned.
func GetCurrentNodeIP() (string,error) {
	// Initialize commands
	//  - `nslookup kubernetes.default.svc.cluster.local` is executed to determine the IP of the current node.
	//  - `nslookup` resolves the DNS name of the Kubernetes service to an IP address.
	//  - The IP address of the current node is the one returned by the DNS lookup.
	cmd := exec.Command("nslookup", "kubernetes.default.svc.cluster.local")
	output, err := cmd.Output()
	if err != nil {
		// If the command fails, return an error.
		return "", fmt.Errorf("nslookup failed: %v", err)
	}

	// Initialize a buffer to store the node IP.
	// The buffer is used to store the output of the `nslookup` command.
	var nodeIpBuffer bytes.Buffer

	// Initialize a scanner to read the output of the `nslookup` command.
	// The scanner is used to iterate over the lines of the output.
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Iterate over the lines of the output.
	// The iteration is used to find the line that contains the IP address of the current node.
	for scanner.Scan() {
		// Get the current line.
		line := scanner.Text()

		// Check if the line contains the "Server" string.
		// If it does, it means that the line contains the IP address of the current node.
		if strings.Contains(line, "Server") {
			// Split the line into fields.
			// The fields are separated by spaces.
			fields := strings.Fields(line)

			// Check if the number of fields is greater than 1.
			// If it is, it means that the line contains the IP address of the current node.
			if len(fields) > 1 {
				// Write the IP address of the current node to the buffer.
				// The IP address is the second field of the line.
				nodeIpBuffer.WriteString(fields[1])
				// Break the loop.
				break
			}
		}
	}

	// Check if the buffer is empty.
	// If it is, it means that the current node's IP was not found.
	if nodeIpBuffer.Len() == 0 {
		// Return an error.
		return "", fmt.Errorf("current node ip not found")
	}

	// Return the NodeIp.
	// The NodeIp is the IP address of the current node.
	return nodeIpBuffer.String(), nil
}
