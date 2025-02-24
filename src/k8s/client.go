package k8s

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient creates a Kubernetes client by loading the kubeconfig from one of the following sources,
// in order of preference:
//
// 1. In-cluster configuration, if running inside a Kubernetes pod.
// 2. The file specified by the KUBECONFIG environment variable.
// 3. The default location, $HOME/.kube/config.
//
// The function returns an error if it fails to create a client.
func GetClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.ExpandEnv("$HOME/.kube/config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to load kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Kubernetes client: %v", err)
	}

	return clientset, nil
}
