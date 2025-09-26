# Troubleshooting

This guide provides solutions to common problems you might encounter while setting up or running the k2ray application.

---

### 1. Backend server fails to start with "address already in use"

*   **Problem**: When you run `go run ./cmd/k2ray`, you see an error message similar to `listen tcp :8080: bind: address already in use`.
*   **Cause**: Another process is already running on the port the backend is trying to use (e.g., port 8080).
*   **Solution**:
    1.  **Identify and stop the other process**: Use a command like `lsof -i :8080` (on macOS/Linux) or `netstat -ano | findstr :8080` (on Windows) to find the Process ID (PID) of the application using the port. Then, stop that application or use `kill <PID>`.
    2.  **Change the port**: If you cannot stop the other process, you can change the port the k2ray backend uses by setting the `PORT` environment variable (Note: This assumes the application is updated to respect the `PORT` variable. As of now, it's hardcoded, but this is a common practice).

---

### 2. Frontend shows a blank page or API errors in the console

*   **Problem**: The frontend loads as a blank white page, or the browser's developer console shows errors like `404 Not Found` for API requests or `CORS error`.
*   **Cause**:
    1.  The backend server is not running.
    2.  The frontend is trying to connect to the wrong backend URL or port.
    3.  A Cross-Origin Resource Sharing (CORS) issue is preventing the frontend from communicating with the backend.
*   **Solution**:
    1.  **Ensure the backend is running**: Make sure you have started the backend server by running `go run ./cmd/k2ray` in a separate terminal.
    2.  **Verify API URL**: Check the frontend code (likely in `web/src/services/api.ts` or a similar file) to ensure the `baseURL` for Axios requests points to the correct backend address (e.g., `http://localhost:8080`).
    3.  **Check CORS Configuration**: Ensure the Gin backend has the correct CORS middleware configuration to allow requests from the frontend's origin (e.g., `http://localhost:5173`).

---

### 3. `npm install` fails in the `web/` directory

*   **Problem**: Running `npm install` inside the `web` directory results in errors.
*   **Cause**: This can be due to an outdated version of Node.js/npm, network issues, or corrupted npm cache.
*   **Solution**:
    1.  **Update Node.js and npm**: Make sure you are using a supported version (v18 or higher).
    2.  **Clear npm cache**: Run `npm cache clean --force` and then try `npm install` again.
    3.  **Check network connection**: Ensure you have a stable internet connection, as npm needs to download packages from the registry.

---

### 4. JWT authentication errors (401 Unauthorized)

*   **Problem**: You are logged out unexpectedly, or API requests fail with a `401 Unauthorized` error.
*   **Cause**:
    1.  The JWT has expired.
    2.  The `JWT_SECRET` used by the backend was changed, invalidating all existing tokens.
    3.  The token is not being sent correctly in the `Authorization` header.
*   **Solution**:
    1.  **Log in again**: The simplest solution is to go back to the login page and re-authenticate to get a new token.
    2.  **Check `JWT_SECRET`**: Ensure the `JWT_SECRET` in your `.env` file or environment variables is correct and has not been accidentally changed.
    3.  **Use browser developer tools**: Check the "Network" tab in your browser's developer tools to inspect the headers of a failing request. Ensure the `Authorization` header is present and has the format `Bearer <your-token>`.