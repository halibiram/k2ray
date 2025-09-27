# Troubleshooting Guide

This guide is intended to help you diagnose and resolve common issues you may encounter while using the application.

## Application Not Starting

- **Symptom:** You run the executable (`k2ray` or `k2ray.exe`) and it closes immediately, or you cannot access the web UI.
- **Solution:**
  1.  **Check the logs:** Run the application from a terminal or command prompt. Any startup errors will be printed there. Look for messages like "Port already in use" or "Failed to open database".
  2.  **Port conflict:** If you see an "address already in use" error, it means another application is using the configured port (default `8080`). Stop the other application or change the `K2RAY_PORT` in your `.env` file.
  3.  **Permissions issue:** Ensure the application has permission to create and write to the database file (`k2ray.db` by default) in its directory. Try running the application as an administrator (on Windows) or with `sudo` (on Linux/macOS) to diagnose permission issues, but be cautious with this approach.

## Cannot Log In

- **Symptom:** You are unable to log in, even with what you believe are the correct credentials.
- **Solution:**
  1.  **Default Credentials:** If this is your first time running the application, use the default credentials: `admin` / `admin`.
  2.  **Caps Lock:** Ensure Caps Lock is not enabled on your keyboard.
  3.  **Check logs:** The application logs might show errors related to database access or JWT secret configuration that could prevent logins.
  4.  **Password Reset:** If you have forgotten your password and there is no password reset feature in the UI, you may need to manually edit the database. This is an advanced procedure and should be done with caution.

## V2Ray Core Issues

- **Symptom:** The dashboard shows that the V2Ray core is not running, or clients cannot connect.
- **Solution:**
  1.  **V2Ray Executable:** Ensure that the V2Ray core executable is correctly placed and configured if the application requires it to be provided manually.
  2.  **Configuration Errors:** A misconfiguration in one of the V2Ray profiles can prevent the core process from starting. Check the application logs for specific error messages from V2Ray.
  3.  **Firewall:** Make sure your server's firewall is not blocking the ports used by your V2Ray configurations.

## Further Assistance

If you are still experiencing issues after following this guide, please consider opening an issue on our [GitHub repository](https://github.com/your-username/k2ray/issues). Please include any relevant logs and a detailed description of the problem.