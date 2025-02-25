package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeInfo struct {
	Name       string
	InternalIP string
}

var currentNodeName string

// GetClusterNodes fetches all nodes in the cluster, filters out the current node by IP and returns a slice of NodeInfo
// containing the name and internal IP of the target nodes. The function returns an error if the Kubernetes client fails to
// list the nodes.
func GetClusterNodes(clientset *kubernetes.Clientset, currentNodeIP string) (string, []NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("Failed to list nodes: %w", err)
	}

	var nodesInfo []NodeInfo
	for _, node := range nodes.Items {
		var internalIP string

		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				internalIP = addr.Address
			}
		}

		if internalIP == currentNodeIP {
			currentNodeName = node.Name
			continue
		}

		nodesInfo = append(nodesInfo, NodeInfo{
			Name:       node.Name,
			InternalIP: internalIP,
		})
	}

	return currentNodeName, nodesInfo, nil
}
