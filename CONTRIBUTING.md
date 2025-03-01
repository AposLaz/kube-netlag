# Contributing to Kube-NetLag

Thank you for your interest in contributing to **Kube-NetLag**! 🎉  
We welcome contributions from everyone, whether it's reporting an issue, improving documentation, or adding new features.

---

## 📌 How to Contribute

We accept **issues**, **pull requests**, and **discussions** to improve Kube-NetLag.  
Here’s how you can contribute:

### **🐛 Reporting Issues**
If you find a **bug**, please:
1. **Check existing issues** to avoid duplicates.
2. **Open a new issue** with:
   - A **clear title** describing the issue.
   - **Steps to reproduce** the bug.
   - Expected vs. actual behavior.
   - Logs or error messages (if applicable).

📌 Open an issue here: [GitHub Issues](https://github.com/AposLaz/kube-netlag/issues)

---

### **🛠️ Submitting a Pull Request (PR)**
1. **Fork the repository** and create a new branch:
   ```sh
   git checkout -b feature/your-feature-name
2. **Make your changes** and commit them with a **descriptive commit message**.
   ```sh
   git commit -m "feat: Add latency histogram support"
   ```
3. **Push your branch** to GitHub:
   ```sh
   git push origin feature/your-feature-name
   ```
4. **Open a Pull Request**
   - Describe your changes.
   - Link related issues (if applicable).
   - Tag relevant maintainers for review.
   
📌 Open a PR here: [GitHub Pull Requests](https://github.com/AposLaz/kube-netlag/pulls)

---

### **💡 Feature Requests**
Want a new feature? Open an issue with:
- A **clear use case** for the feature.
- Any **alternatives** you considered.
- How it aligns with Kube-NetLag’s goals.

---

### **🏗️ Code Guidelines**
To ensure high-quality contributions, follow these best practices:
#### **Code Style**
- Follow **Go coding conventions**.
- Use **meaningful variable names**.
- Format your code with:
  ```sh
  go fmt ./..
  ```
---

#### **Project Structure**
- **Put Kubernetes-related code in** `/k8s`.
- **Prometheus metrics go in** `/promMetrics`.
- **Network performance logic stays in** `/netperf`.

---

### **🛠️ Development Setup**
To set up Kube-NetLag for development:
1. **Clone the repository**:
   ```sh
   git clone https://github.com/AposLaz/kube-netlag.git
   cd kube-netlag
   ```
2. **Intall dependencies**:
   ```sh
   go mod tidy
   ```
3. **Run the application**:
   ```sh
   go run .
   ```
---

## 📜 Code of Conduct
We follow the
Be respectful and inclusive to all community members.

---

## License
By contributing, you agree that your code will be licensed under the Apache 2.0 License.
Read the full license: [LICENSE](./LICENSE)

## Need Help?
- **Ask a question** in [Discussions](https://github.com/AposLaz/kube-netlag/discussions).
- **Report a security issue** via [GitHub Security Advisories](https://github.com/AposLaz/kube-netlag/security/advisories).
