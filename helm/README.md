# Kube-NetLag

**Kube-NetLag** is a lightweight network performance testing tool deployed as a **DaemonSet** across your Kubernetes cluster. It measures **inter-node network latency** using [Netperf](https://github.com/HewlettPackard/netperf) and exposes results as **Prometheus metrics**, making it easy to monitor and alert on cluster network health.

- **Source Code:** [github.com/AposLaz/kube-netlag](https://github.com/AposLaz/kube-netlag)
- **Artifact Hub:** [artifacthub.io/packages/helm/kube-netlag/kube-netlag](https://artifacthub.io/packages/helm/kube-netlag/kube-netlag)

---

## Features

- Measures **min / max / average** latency between every pair of nodes.
- Exposes **Prometheus metrics** at `:9090/metrics` on each node.
- Optional **automatic Prometheus ConfigMap** generation for scrape configuration.
- Fully configurable via `values.yaml`.
- Runs as a non-root, read-only container (security-hardened by default).

---

## Prerequisites

| Requirement | Version |
|---|---|
| Kubernetes | 1.26+ |
| Helm | 3.0+ |
| Prometheus | Any (optional, for metrics) |

---

## Installation

### 1. Add the Helm repository

```sh
helm repo add kube-netlag https://aposlaz.github.io/kube-netlag
helm repo update
```

### 2. Install the chart

Install with **default values** into a dedicated namespace:

```sh
helm install kube-netlag kube-netlag/kube-netlag \
  --namespace kube-netlag \
  --create-namespace
```

Install with a **custom values file:**

```sh
helm install kube-netlag kube-netlag/kube-netlag \
  --namespace kube-netlag \
  --create-namespace \
  -f my-values.yaml
```

Install and **place the Prometheus ConfigMap in a specific namespace** (e.g. where Prometheus lives):

```sh
helm install kube-netlag kube-netlag/kube-netlag \
  --namespace kube-netlag \
  --create-namespace \
  --set prometheusConfig.create=true \
  --set prometheusConfig.namespace=monitoring
```

---

## Upgrading

```sh
helm repo update
helm upgrade kube-netlag kube-netlag/kube-netlag --namespace kube-netlag
```

---

## Uninstalling

```sh
helm uninstall kube-netlag --namespace kube-netlag
kubectl delete namespace kube-netlag
```

---

## Verifying the Deployment

**Check that all DaemonSet pods are running:**

```sh
kubectl get pods -n kube-netlag -l app.kubernetes.io/name=kube-netlag
```

**Tail logs to confirm measurements are running:**

```sh
kubectl logs -l app.kubernetes.io/name=kube-netlag -n kube-netlag --follow
```

**Manually scrape metrics from a pod:**

```sh
kubectl port-forward -n kube-netlag daemonset/kube-netlag 9090:9090
curl http://localhost:9090/metrics
```

---

## Configuration

The following table lists all configurable parameters with their defaults.

### Global

| Parameter | Description | Default |
|---|---|---|
| `image.repository` | Container image repository | `alazidis/kube-netlag` |
| `image.tag` | Image tag (defaults to chart `appVersion`) | `""` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `imagePullSecrets` | List of image pull secret names | `[]` |
| `nameOverride` | Override the chart name | `""` |
| `fullnameOverride` | Override the full release name | `""` |
| `namespaceOverride` | Override the deployment namespace | `""` |

### Ports

| Parameter | Description | Default |
|---|---|---|
| `ports[0].containerPort` | Netperf server port (container) | `12865` |
| `ports[0].hostPort` | Netperf server port (host) | `12865` |
| `ports[1].containerPort` | Prometheus metrics port (container) | `9090` |
| `ports[1].hostPort` | Prometheus metrics port (host) | `9090` |

> **Note:** If you change either port, add the corresponding `NETPERF_PORT` or `METRICS_PORT` entry in `extraEnv`.

### Extra Environment Variables

| Parameter | Description | Default |
|---|---|---|
| `extraEnv` | Additional environment variables injected into the container | `{}` |

**Example — override ports via env vars:**

```yaml
extraEnv:
  - name: NETPERF_PORT
    value: "12865"
  - name: METRICS_PORT
    value: "9090"
```

### Prometheus Configuration

| Parameter | Description | Default |
|---|---|---|
| `prometheusConfig.create` | Create a Prometheus scrape ConfigMap | `true` |
| `prometheusConfig.namespace` | Namespace to deploy the ConfigMap into (defaults to release namespace) | `""` |

When `prometheusConfig.create: true`, a ConfigMap is created that configures Prometheus to automatically discover and scrape kube-netlag pods via Kubernetes service discovery.

### Service Account

| Parameter | Description | Default |
|---|---|---|
| `serviceAccount.annotations` | Annotations to add to the ServiceAccount | `{}` |

### Pod Configuration

| Parameter | Description | Default |
|---|---|---|
| `podAnnotations` | Annotations added to each pod | See below |
| `podLabels` | Extra labels added to each pod | `{}` |

Default `podAnnotations`:

```yaml
podAnnotations:
  prometheus.io/path: /metrics
  prometheus.io/port: "9090"
  prometheus.io/scrape: "true"
```

### Security Context

| Parameter | Description | Default |
|---|---|---|
| `podSecurityContext.fsGroup` | Volume GID | `1000` |
| `podSecurityContext.runAsUser` | UID to run as | `1000` |
| `podSecurityContext.runAsGroup` | GID to run as | `1000` |
| `podSecurityContext.runAsNonRoot` | Enforce non-root execution | `true` |
| `securityContext.capabilities.drop` | Dropped Linux capabilities | `["ALL"]` |
| `securityContext.readOnlyRootFilesystem` | Read-only root filesystem | `true` |
| `securityContext.allowPrivilegeEscalation` | Prevent privilege escalation | `false` |
| `securityContext.runAsNonRoot` | Enforce non-root in container | `true` |

### Resources

| Parameter | Description | Default |
|---|---|---|
| `resources.requests.cpu` | CPU request | `200m` |
| `resources.requests.memory` | Memory request | `40Mi` |
| `resources.limits.cpu` | CPU limit | `400m` |
| `resources.limits.memory` | Memory limit | `80Mi` |

### Probes

| Parameter | Description | Default |
|---|---|---|
| `livenessProbe.initialDelaySeconds` | Liveness probe initial delay | `10` |
| `livenessProbe.periodSeconds` | Liveness probe interval | `10` |
| `readinessProbe.initialDelaySeconds` | Readiness probe initial delay | `5` |
| `readinessProbe.periodSeconds` | Readiness probe interval | `5` |

---

## Prometheus Integration

When `prometheusConfig.create: true`, the chart creates a ConfigMap containing a ready-to-use Prometheus scrape configuration that uses **Kubernetes pod service discovery** to automatically find all kube-netlag pods:

```yaml
scrape_configs:
  - job_name: "kube-netlag-daemon"
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - <release-namespace>
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_name]
        action: keep
        regex: kube-netlag
      ...
```

Mount or reference this ConfigMap from your Prometheus deployment to enable scraping.

---

## Exposed Prometheus Metrics

All metrics include node-pair labels for fine-grained visibility.

| Metric | Description | Unit |
|---|---|---|
| `node_min_latency_ms` | Minimum latency between two nodes | microseconds |
| `node_max_latency_ms` | Maximum latency between two nodes | microseconds |
| `node_avg_latency_ms` | Average latency between two nodes | microseconds |

**Labels on every metric:**

| Label | Description |
|---|---|
| `from_node` | Name of the source node |
| `to_node` | Name of the destination node |
| `from_ip` | IP address of the source node |
| `to_ip` | IP address of the destination node |

**Example PromQL queries:**

```promql
# Average latency from node-1 to node-2
node_avg_latency_ms{from_node="node-1", to_node="node-2"}

# Top 5 highest average latencies across the cluster
topk(5, node_avg_latency_ms)

# Alert: any node pair with avg latency above 5ms
node_avg_latency_ms > 5
```

---

## Tested Kubernetes Distributions

| Distribution | Status |
|---|---|
| Rancher Kubernetes | ✅ Tested |
| AWS EKS            | ✅ Tested |
| Others | Community contributions welcome |

---

## Contributing & Support

- [Contribution Guidelines](https://github.com/AposLaz/kube-netlag/blob/main/CONTRIBUTING.md)
- [Open an Issue](https://github.com/AposLaz/kube-netlag/issues)
- [Submit a Pull Request](https://github.com/AposLaz/kube-netlag/pulls)

---

## License

Kube-NetLag is licensed under the **Apache 2.0 License**.

---

## Maintainer

**Apostolos Lazidis** - [aplazidis@gmail.com](mailto:aplazidis@gmail.com)
