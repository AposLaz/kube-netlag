/*
Copyright 2024 Apostolos Lazidis

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"context"
	"fmt"
	"strings"

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

		if strings.Contains(node.Name, "master") || strings.Contains(node.Name, "control-plane") {
			continue
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
