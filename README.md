# **Kube-NetLag**

<p align="center">
    <img src="kube-netlag.png" alt="Kube-NetLag Logo" width="300"/>
</p>

# **Table of Contents**

- [**Kube-NetLag**](#kube-netlag)
  - [**Overview**](#overview)
  - [**Features**](#features)
  - [**Tested Kubernetes Distributions**](#tested-kubernetes-distributions)
  - [**Installation**](#installation)
    - [**Prerequisites**](#prerequisites)
    - [**Deploy using Kubernetes Manifests**](#deploy-using-kubernetes-manifests)
    - [**Deploy using Helm Chart**](#deploy-using-helm-chart)
      - [**Install the Chart**](#install-the-chart)
  - [**Testing the Deployment**](#testing-the-deployment)
    - [**Verify that the DaemonSet is running**](#verify-that-the-daemonset-is-running)
    - [**Check logs to confirm proper operation**](#check-logs-to-confirm-proper-operation)
  - [**Uninstall the Chart**](#uninstall-the-chart)
  - [**Configuration**](#configuration)
    - [**Global Parameters**](#global-parameters)
    - [**Network & Ports**](#network--ports)
    - [**Prometheus Configuration**](#prometheus-configuration)
    - [**Resource Allocation**](#resource-allocation)
  - [**Prometheus Integration**](#prometheus-integration)
  - [**Exposed Prometheus Metrics**](#exposed-prometheus-metrics)
    - [**Latency Metrics**](#latency-metrics)
    - [**Example Prometheus Query**](#example-prometheus-query)
  - [**Contributing**](#contributing)
  - [**License**](#license)
  - [**Author**](#author)

## **Overview**

Kube-NetLag is a lightweight **network performance testing** and **latency simulation** tool deployed as a DaemonSet across your Kubernetes cluster. It helps measure **network delays between nodes**, making it useful for **debugging and optimizing cluster communication**.  

You can deploy **Kube-NetLag** using either Kubernetes **manifests** or a **Helm chart**. The Helm chart simplifies deployment and configuration, enabling seamless monitoring with **Prometheus integration**.  

## **Features**
- ✅ **Network performance testing** using Netperf.
- 📊 **Prometheus metrics** exposure for monitoring.
- ⚙️ **Customizable ports & environment variables** for flexibility.
- 📡 **Automatic Prometheus configuration** (optional).

## **Tested Kubernetes Distributions**
Kube-NetLag has been successfully tested on the following Kubernetes distributions:
- ✅ **Rancher Kubernetes**
- 🔄 **More distributions to be tested...**

> If you tested it on another Kubernetes distribution, feel free to contribute and share your feedback!

## **Installation**

### **Prerequisites**
- Kubernetes **1.26+**
- Helm **3.0+**
- Prometheus (if monitoring is required)

### Deploy using Kubernetes Manifests

The manifest files are located in the `manifests` directory.

```sh
kubectl apply -f manifests/
```
### Deploy using Helm Chart

#### **Install the Chart**

To install Kube-NetLag with default values:

```sh
helm repo add alazidis https://your-helm-repo-url
helm install kube-netlag alazidis/kube-netlag
```

To install with a custom configuration:

```sh
helm install kube-netlag alazidis/kube-netlag -f values.yaml
```

### **Testing the Deployment**

#### **Verify that the DaemonSet is running:**
```sh
kubectl get pods -n kube-netlag
```
#### Check logs to confirm proper operation:

```sh
kubectl logs -l app.kubernetes.io/name=kube-netlag -n kube-netlag
```
### Uninstall the Chart

To completely remove Kube-NetLag, run:

```sh
helm uninstall kube-netlag
kubectl delete namespace kube-netlag
```
### **Configuration**
The Helm chart allows full customization via the `values.yaml` file.

#### **Global Parameters**
| Parameter              | Description                              | Default                     |
|------------------------|------------------------------------------|-----------------------------|
| `image.repository`     | Docker image repository                 | `alazidis/kube-netlag`      |
| `image.tag`           | Image tag                                | `Chart.Version`                    |
| `image.pullPolicy`    | Image pull policy                        | `Always`                    |
| `namespaceOverride`   | Overrides the namespace for deployment   | `""`                        |

---

#### **Network & Ports**
| Parameter              | Description                            | Default  |
|------------------------|----------------------------------------|----------|
| `ports.containerPort`  | Netperf service port                  | `12865`  |
| `ports.hostPort`       | Host-mapped Netperf port              | `12865`  |
| `ports.containerPort`  | Prometheus metrics port               | `9090`   |
| `ports.hostPort`       | Host-mapped Prometheus port           | `9090`   |

> **Note:** If modifying ports, update `extraEnv` accordingly.

---

#### **Prometheus Configuration**

| Parameter                 | Description                               | Default  |
|---------------------------|-------------------------------------------|----------|
| `prometheusConfig.create` | Enable automatic Prometheus ConfigMap    | `true`   |
| `prometheusConfig.namespace` | Namespace for Prometheus ConfigMap  | `""`     |

---

#### **Resource Allocation**
| Parameter                   | Description            | Default  |
|-----------------------------|------------------------|----------|
| `resources.requests.cpu`    | Requested CPU         | `200m`   |
| `resources.requests.memory` | Requested Memory      | `40Mi`   |
| `resources.limits.cpu`      | CPU limit             | `400m`   |
| `resources.limits.memory`   | Memory limit          | `80Mi`   |

---

#### **Prometheus Integration**

If `prometheusConfig.create` is set to `true`, a Prometheus scrape job will be automatically created:

```yaml
- job_name: "kube-netlag-daemon"
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - {{ .Values.namespaceOverride | default .Release.Namespace }}
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_app]
      action: keep
      regex: {{ include "..name" . }}
    - source_labels: [__meta_kubernetes_namespace]
      action: keep
      regex: {{ .Values.namespaceOverride | default .Release.Namespace }}
      .....
```

Prometheus will automatically collect network latency and performance metrics from Kube-NetLag.

## **Exposed Prometheus Metrics**

Kube-NetLag provides the following **Prometheus metrics** to monitor network latency between Kubernetes nodes.

### **Latency Metrics**
| Metric Name               | Description                                           | Labels (`from_node`, `to_node`, `from_ip`, `to_ip`) |
|---------------------------|------------------------------------------------------|------------------------------------------------------|
| `node_min_latency_ms`     | Minimum latency in **microseconds** between nodes.  | ✅ |
| `node_max_latency_ms`     | Maximum latency in **microseconds** between nodes.  | ✅ |
| `node_avg_latency_ms`     | Average latency in **microseconds** between nodes.  | ✅ |

Each metric includes the following labels:
- **`from_node`** – Name of the source node (The current Node).
- **`to_node`** – Name of the destination node.
- **`from_ip`** – IP address of the source node.
- **`to_ip`** – IP address of the destination node.

### **Example Prometheus Query**
To visualize average latency between nodes in Prometheus:

```promql
node_avg_latency_ms{from_node="node-1", to_node="node-2"}
```

# Contributing

Contributions are welcome! To report issues or request features:

- Open an issue on GitHub.
- Submit a pull request with improvements.

# License

This project is licensed under the MIT License.

# Author

- Apostolos Lazidis 🚀