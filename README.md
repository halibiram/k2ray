# k2ray

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/halibiram/k2ray)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue)](https://github.com/halibiram/k2ray/releases)

k2ray is a modern, web-based management panel for the V2Ray proxy tool. It provides a user-friendly interface to manage your V2Ray configurations, monitor system status, and view traffic metrics, all from your web browser.

---

## 📸 Screenshots

*(Add screenshots or GIFs of the application here to give users a visual preview.)*

![K2Ray Screenshot](https://place-hold.it/800x450/663399/ffffff?text=k2ray%20UI%20Screenshot)

---

## ✨ Features

*   **Advanced Configuration Management**: Easily create, edit, delete, and switch between multiple V2Ray configurations.
*   **QR Code Integration**: Seamlessly export and import configurations using QR codes. Scan with your mobile device's camera or upload an image file.
*   **Real-Time System Monitoring**: Keep an eye on system status, including CPU, memory, and disk usage.
*   **Live Log Viewer**: View and search V2Ray logs directly from the web interface.
*   **Traffic Metrics & Analytics**: Monitor real-time network traffic, connection statistics, and data usage with interactive charts.
*   **Multi-Language Support**: Fully localized interface with support for both English and Turkish.
*   **Customizable Themes**: Switch between light and dark themes for a comfortable user experience.
*   **Secure Access**: Protect your panel with JWT-based authentication and optional Two-Factor Authentication (2FA).
*   **RESTful API**: A well-documented API for programmatic access and integration.

---

## 🛠️ Getting Started

This section provides a brief overview of how to get the project running. For more detailed instructions, please refer to the documentation in the `docs/` directory.

### ✅ System Requirements

*   **Go**: Version 1.24 or higher
*   **Node.js**: Version 18 or higher (with npm)
*   **V2Ray**: A running instance of V2Ray that k2ray can manage.

### ⚙️ Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/halibiram/k2ray.git
    cd k2ray
    ```

2.  **Setup the Backend (API Server):**
    ```bash
    # Run the development server
    go run ./cmd/k2ray
    ```
    The API server will start, typically on port `8080`.

3.  **Setup the Frontend (Web UI):**
    ```bash
    # Navigate to the web directory
    cd web

    # Install dependencies
    npm install

    # Run the development server
    npm run dev
    ```
    The frontend development server will be accessible at `http://localhost:5173`.

---

## 📈 Monitoring

For a comprehensive guide on setting up logging, metrics, and alerting, please see the [Monitoring Documentation](./docs/monitoring.md).

---

## 🧪 Testing

This project is equipped with a comprehensive suite of tests and quality assurance tools to ensure code reliability and maintainability.

### Backend (Go)

All backend tests can be run from the root of the project using the `Makefile`.

*   **Run all unit and integration tests:**
    ```bash
    make test
    ```

*   **Run tests with code coverage:**
    This will generate a `coverage.out` file, which can be used by code analysis tools.
    ```bash
    make test-coverage
    ```

*   **Run the linter (`golangci-lint`):**
    To generate a JSON report for SonarQube integration, use:
    ```bash
    make lint-report
    ```

### Frontend (Vue.js)

All frontend tests and scripts should be run from within the `web/` directory.

*   **Run component tests (`vitest`):**
    ```bash
    cd web
    npm test
    ```

*   **Run component tests with code coverage:**
    This generates reports in `web/tests/coverage/`, including `lcov.info` for SonarQube.
    ```bash
    cd web
    npm run test:coverage
    ```

*   **Run the E2E Smoke Test:**
    This test verifies that the frontend and backend servers can start and that the frontend serves the correct initial page.
    *Note: This requires both the backend and frontend servers to be running.*
    ```bash
    cd web
    npm run test:smoke
    ```

*   **Lint and Format Code:**
    To check for linting errors:
    ```bash
    cd web
    npm run lint
    ```
    To automatically format the code with Prettier:
    ```bash
    cd web
    npm run format
    ```
    To generate a JSON report for SonarQube integration:
    ```bash
    cd web
    npm run lint:report
    ```

---

## 🤝 Contribution Guidelines

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

---

## 📄 License

This project is distributed under the MIT License. See `LICENSE` for more information.

---

## 📜 Changelog

For a detailed list of changes, please see the [CHANGELOG.md](./CHANGELOG.md) file.