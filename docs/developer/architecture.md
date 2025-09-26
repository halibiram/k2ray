# Architecture

This document provides a high-level overview of the k2ray project's architecture. The project is a modern web application with a Go backend and a Vue.js frontend.

## 1. High-Level Overview

The system is composed of two main parts:

1.  **Backend API Server**: A Go application responsible for business logic, data persistence, user authentication, and managing the V2Ray process.
2.  **Frontend Web Client**: A Vue.js Single-Page Application (SPA) that provides the user interface for interacting with the backend.

These two components are decoupled and communicate via a RESTful JSON API.

![High-Level Architecture Diagram](https://i.imgur.com/example.png) <!-- Placeholder for a real diagram -->

---

## 2. Backend Architecture

The backend is built using **Go** and follows a standard project layout for Go applications.

*   **Web Framework**: **Gin** (`github.com/gin-gonic/gin`) is used as the HTTP web framework. It provides fast routing, middleware support, and request/response handling. API routes are defined in `internal/api/routes.go`.
*   **Database**: **SQLite** (`github.com/mattn/go-sqlite3`) is used as the database for simplicity and ease of setup. The database connection and query logic are located in the `internal/db` package.
*   **Authentication**: User authentication is handled using **JSON Web Tokens (JWT)**. The `internal/auth` package contains the logic for generating, parsing, and validating tokens. Middleware in `internal/api/middleware` protects authenticated routes.
*   **Configuration**: Application configuration is managed via environment variables and a `.env` file (see `internal/config/manager.go`). This allows for flexible configuration in different environments.
*   **V2Ray Integration**: The backend directly manages the V2Ray core process. The logic for starting, stopping, and monitoring the process is located in `internal/v2ray`.

### Key Backend Directories:

*   `cmd/k2ray/`: The main entry point of the application.
*   `internal/`: Contains all the core application logic.
    *   `api/`: Defines API routes and handlers.
    *   `auth/`: JWT authentication logic.
    *   `config/`: Application configuration management.
    *   `db/`: Database initialization and interaction.
    *   `v2ray/`: V2Ray process management.
*   `configs/`: Default location for environment configuration files.

---

## 3. Frontend Architecture

The frontend is a modern Single-Page Application (SPA) built with **Vue.js**.

*   **Framework**: **Vue 3** with the Composition API.
*   **Build Tool**: **Vite** is used for fast development server startup and optimized production builds. Configuration is in `vite.config.ts`.
*   **Routing**: **Vue Router** (`vue-router`) is used for client-side routing. Route definitions are in `web/src/router/index.ts`.
*   **State Management**: **Pinia** (`pinia`) is used for centralized state management. Stores are defined in `web/src/stores/`.
*   **Styling**: **Tailwind CSS** is used for utility-first styling. The configuration is in `web/tailwind.config.js`.
*   **API Communication**: **Axios** is used for making HTTP requests to the backend API. A pre-configured Axios instance is typically set up in a service module (e.g., `web/src/services/api.ts`).

### Key Frontend Directories (`web/src`):

*   `assets/`: Static assets like images and global CSS.
*   `components/`: Reusable Vue components (e.g., buttons, modals, input fields).
*   `views/`: Page-level components that are mapped to routes (e.g., `HomeView.vue`, `LoginView.vue`).
*   `router/`: Vue Router configuration.
*   `stores/`: Pinia state management stores.
*   `services/`: Modules for interacting with external services, like the backend API.
*   `main.ts`: The entry point of the Vue application.