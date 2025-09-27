# Frequently Asked Questions (FAQ)

This page answers some of the most common questions about the K2Ray project.

---

### **Q: What is K2Ray?**

**A:** K2Ray is a web-based management panel for V2Ray. It provides a user-friendly interface to help you create, manage, and monitor V2Ray configurations and users, simplifying the process of running a V2Ray server.

---

### **Q: Is K2Ray free?**

**A:** Yes, K2Ray is open-source software and is free to use. You can find the source code on its GitHub repository.

---

### **Q: What technologies does K2Ray use?**

**A:** The application is built with a Go backend and a Vue.js frontend. It uses SQLite as its default database for simplicity.

---

### **Q: How do I update K2Ray to a new version?**

**A:** To update, you can download the latest release from the [GitHub Releases page](https://github.com/your-username/k2ray/releases), stop the currently running application, replace the old executable with the new one, and then restart it. Your database (`k2ray.db`) and `.env` file will be preserved.

---

### **Q: How do I report a bug or suggest a feature?**

**A:** We appreciate your feedback! You can report bugs or suggest new features by creating an "Issue" on our [GitHub repository issues page](https://github.com/your-username/k2ray/issues). Please provide as much detail as possible.

---

### **Q: Where is my data stored?**

**A:** By default, all your data (users, configurations, etc.) is stored in a single file named `k2ray.db` in the same directory as the application. You can change the path to this file using the `K2RAY_DB_PATH` environment variable.

---

### **Q: How can I contribute to the project?**

**A:** We welcome contributions! If you are a developer, please read our [Contributing Guide](../developer/contributing.md) to get started.