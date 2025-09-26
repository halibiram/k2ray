# k2ray

k2ray is a modern, web-based management panel for the V2Ray proxy tool. It provides a user-friendly interface to manage your V2Ray configurations, monitor system status, and view traffic metrics, all from your web browser.

![K2Ray Screenshot](https://place-hold.it/800x450/663399/ffffff?text=k2ray%20UI%20Screenshot)

---

## ✨ Key Features

*   **Configuration Management**: Easily create, edit, delete, and switch between V2Ray configurations.
*   **QR Code Support**: Export and import configurations seamlessly using QR codes—scan with your camera or upload an image file.
*   **System Monitoring**: View real-time system status, logs, and resource usage.
*   **Traffic Metrics**: Monitor network traffic and connection statistics.
*   **Multi-Language Support**: Switch between English and Turkish interfaces.
*   **Light & Dark Themes**: Choose your preferred visual theme for a comfortable user experience.

---

## 🚀 Technology Stack

*   **Backend**: Go with the Gin web framework.
*   **Frontend**: Vue.js 3 with Vite, Pinia for state management, and Tailwind CSS for styling.
*   **Database**: SQLite for simple, file-based data storage.
*   **API Documentation**: OpenAPI 3.0 specification.

---

## 🛠️ Getting Started

This section provides a brief overview of how to get the project running. For more detailed instructions, please refer to the documentation in the `docs/` directory.

### Prerequisites

*   Go (version 1.21+)
*   Node.js (version 18+) and npm

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/k2ray.git
    cd k2ray
    ```

2.  **Setup the Backend:**
    ```bash
    # Run the development server
    go run ./cmd/k2ray
    ```

3.  **Setup the Frontend:**
    ```bash
    # Navigate to the web directory
    cd web

    # Install dependencies
    npm install

    # Run the development server
    npm run dev
    ```

After these steps, the backend API will be running (typically on port 8080) and the frontend will be accessible at `http://localhost:5173`.

---

## 📚 Documentation

This project includes comprehensive documentation to help you get started:

*   **[Installation Guide](./docs/user/installation.md)**: Step-by-step instructions for setting up the project.
*   **[Configuration Guide](./docs/user/configuration.md)**: How to configure the application.
*   **[API Reference](./docs/api/openapi.yaml)**: Full OpenAPI specification for the backend API.
*   **[Architecture Overview](./docs/developer/architecture.md)**: A look into the project's structure for developers.
*   **[Troubleshooting](./docs/user/troubleshooting.md)**: Solutions for common issues.
*   **[FAQ](./docs/user/faq.md)**: Frequently asked questions.

---

## 🤝 Contributing

Contributions are welcome! If you'd like to contribute, please fork the repository and create a pull request. For major changes, please open an issue first to discuss what you would like to change.