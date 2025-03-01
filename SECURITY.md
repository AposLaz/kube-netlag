# 🛡️ Security Policy

## 🔹 Supported Versions
The following table shows the versions of Kube-NetLag currently receiving **security updates**.

| Version  | Supported          |
|----------|--------------------|
| `latest` | ✅ Actively supported |
| `0.1.x`  | ✅ Security patches only |
| `< 0.1`  | ❌ No longer supported |

If you are using an **unsupported version**, we strongly recommend upgrading to the **latest release**.

---

## 🚨 Reporting a Vulnerability

We take security issues **very seriously**. If you discover a vulnerability in **Kube-NetLag**, please follow these steps:

### 📩 **How to Report**
1. **DO NOT** create a public GitHub issue for security vulnerabilities.
2. Instead, **email** us at **[aplazidis@gmail.com](mailto:aplazidis@gmail.com)** with:
   - A detailed description of the vulnerability.
   - Steps to reproduce the issue.
   - Potential impact and severity assessment.
   - Any suggested fixes (if available).

### 🔒 **Responsible Disclosure**
- We will **acknowledge your report within 48 hours**.
- A fix will be developed **privately** and released in a security patch.
- You will be **credited** in the release notes (unless you wish to remain anonymous).
- If the issue is **critical**, we may **coordinate disclosure** with the CNCF or Kubernetes security teams.

---

## ✅ **Security Best Practices**
To keep your **Kube-NetLag** deployment secure:
- **Use the latest version** (check [releases](https://github.com/AposLaz/kube-netlag/releases)).
- **Follow the principle of least privilege** for Kubernetes RBAC.
- **Monitor Prometheus metrics** for unexpected network behavior.
- **Use TLS encryption** for secure communication (if applicable).
- **Regularly update your Kubernetes cluster**.

---

## 🛠️ **Security Tools**
We encourage users to test Kube-NetLag with **security tools** like:
- [Trivy](https://github.com/aquasecurity/trivy) – Container security scanning.
- [Falco](https://github.com/falcosecurity/falco) – Runtime security monitoring.

---

If you have any **security concerns** or suggestions, feel free to reach out! 🚀
