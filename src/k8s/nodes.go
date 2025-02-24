package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeInfo struct {
	Name            string
	InternalIP      string
	CurrentNodeName string
}

func GetClusterNodes(clientset *kubernetes.Clientset, currentNodeIP string) ([]NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list nodes: %w", err)
	}

	var nodesInfo []NodeInfo
	for _, node := range nodes.Items {
		var internalIP string
		var currentNodeName string

		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				internalIP = addr.Address
				currentNodeName = node.Name
			}
		}

		if internalIP == currentNodeIP {
			continue
		}

		nodesInfo = append(nodesInfo, NodeInfo{
			Name:            node.Name,
			InternalIP:      internalIP,
			CurrentNodeName: currentNodeName,
		})
	}

	return nodesInfo, nil
}
