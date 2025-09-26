# Installation

This guide will walk you through the process of setting up the k2ray project for development and production. The project consists of a Go backend and a Vue.js frontend.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Go**: Version 1.21 or higher. You can download it from [golang.org](https://golang.org/).
*   **Node.js and npm**: Version 18 or higher. You can download it from [nodejs.org](https://nodejs.org/).
*   **Git**: For cloning the repository.

## 1. Clone the Repository

First, clone the project repository from GitHub:

```bash
git clone <repository-url>
cd k2ray
```

## 2. Backend Setup (Go)

The backend is a Go application that serves the API.

### Running in Development

To run the backend server in development mode, navigate to the root of the project and use the `go run` command:

```bash
go run ./cmd/k2ray
```

This will start the backend server, typically on a port defined in the application's configuration (e.g., `8080`).

### Building for Production

To create a production build, use the `go build` command. This will compile the application into a single executable file.

```bash
go build -o k2ray_server ./cmd/k2ray
```

After the build is complete, you can run the application with:

```bash
./k2ray_server
```

## 3. Frontend Setup (Vue.js)

The frontend is a Vue.js single-page application located in the `web/` directory.

### Install Dependencies

Navigate to the `web` directory and install the required npm packages:

```bash
cd web
npm install
```

### Running in Development

To start the frontend development server with hot-reloading, run the following command:

```bash
npm run dev
```

This will typically make the frontend available at `http://localhost:5173`. The development server will automatically proxy API requests to the backend server.

### Building for Production

To build the frontend assets for production, run:

```bash
npm run build
```

This command will create a `dist/` directory inside `web/` containing the optimized static files (HTML, CSS, JavaScript). These files are ready to be served by a web server like Nginx or by the Go backend itself.