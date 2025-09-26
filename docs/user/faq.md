# FAQ (Frequently Asked Questions)

Here are some frequently asked questions about the k2ray project.

---

### **Q1: What is k2ray?**

**A1:** k2ray is a web-based management panel for V2Ray, a powerful network proxy tool. It provides a user-friendly interface to manage V2Ray configurations, monitor system status, and view traffic metrics, all from a web browser.

---

### **Q2: What technologies is k2ray built with?**

**A2:** The project is built with a modern technology stack:
*   **Backend**: Go (with the Gin framework)
*   **Frontend**: Vue.js 3 (with Vite, Pinia, and Tailwind CSS)
*   **Database**: SQLite (for simplicity and portability)

---

### **Q3: Can I use a different database like PostgreSQL or MySQL?**

**A3:** Currently, the application is designed to work with SQLite out of the box. While the database logic is encapsulated in the `internal/db` package, adding support for other databases would require code modifications to include the appropriate database driver and potentially adjust SQL queries.

---

### **Q4: Is the application secure?**

**A4:** The application uses JWT (JSON Web Tokens) for authenticating users, which is a standard and secure method for APIs. However, the overall security of your installation depends on your environment and configuration. **It is critical to change the default `JWT_SECRET`** to a strong, unique value for any production deployment.

---

### **Q5: How do I run k2ray in production?**

**A5:** For production, you should:
1.  **Build the Go backend**: Run `go build -o k2ray_server ./cmd/k2ray` to create a compiled executable.
2.  **Build the frontend assets**: Run `npm run build` in the `web/` directory.
3.  **Serve the frontend**: Configure the Go backend to serve the static files from the `web/dist` directory, or use a dedicated web server like Nginx or Caddy to serve the frontend and proxy API requests to the backend server.
4.  **Use a process manager**: Run the backend executable (`k2ray_server`) using a process manager like `systemd` or `supervisor` to ensure it runs continuously and restarts automatically if it crashes.
5.  **Set environment variables**: Do not rely on the `.env` file in production. Set the `DATABASE_URL` and `JWT_SECRET` as environment variables in your deployment environment.

---

### **Q6: How can I contribute to the project?**

**A6:** Contributions are welcome! Please refer to the `CONTRIBUTING.md` file (if available) for guidelines on how to contribute. Typically, you would fork the repository, create a new branch for your feature or bug fix, and then submit a pull request.