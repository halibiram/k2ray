# Development Setup Guide

This guide describes how to set up a local development environment for this project to contribute code.

## Prerequisites

Before you begin, ensure you have the following installed on your system:
- **Git:** For version control.
- **Go:** v1.18 or later.
- **Node.js:** v18.x or later (for the frontend).
- **npm:** v9.x or later (for the frontend).
- **A code editor:** Such as Visual Studio Code with the Go extension.

## 1. Fork and Clone the Repository

First, fork the main repository on GitHub, and then clone your fork to your local machine:

```bash
git clone https://github.com/your-username/k2ray.git
cd k2ray
```

## 2. Backend Setup (Go)

The backend is a Go application.

### a. Configure Environment Variables

The backend uses a `.env` file for configuration.

1.  Navigate to the `configs/` directory and copy the example file:
    ```bash
    cp .env.example .env
    ```
2.  Open the `.env` file and review the settings. The default settings are generally suitable for local development.

### b. Install Dependencies & Run

Go modules are used to manage dependencies. They will be downloaded automatically when you build or run the application.

1.  Navigate to the root of the project.
2.  Build and run the backend server:
    ```bash
    go run ./cmd/k2ray
    ```
3.  The backend API should now be running, typically on port 8080 as specified in the `.env` file.

## 3. Frontend Setup (Vue.js)

The frontend is a Vue.js application located in the `web/` directory.

1.  Navigate to the frontend directory:
    ```bash
    cd web
    ```
2.  Install the JavaScript dependencies:
    ```bash
    npm install
    ```
3.  Start the frontend development server (powered by Vite):
    ```bash
    npm run dev
    ```
4.  The frontend application should now be running, typically on `http://localhost:5173`, and it will proxy API requests to the backend server running on port 8080.

## 4. Running Tests

### Backend Tests

To run the Go tests for the backend:
```bash
go test ./...
```

### Frontend Tests

To run the frontend tests (if any are configured):
```bash
cd web
npm test
```

You are now ready to start developing! Remember to create a new branch for your changes.