# Security Policy

The security of k2ray is a top priority. We are committed to ensuring our application is secure and our users' data is protected. This document outlines our security features, our policy for reporting vulnerabilities, and best practices for developers.

## Implemented Security Features

This application includes a multi-layered security approach to protect against common web application threats.

### 1. Authentication and Access Control

*   **JWT-Based Authentication:** The application uses JSON Web Tokens (JWT) for secure, stateless authentication. Access tokens are short-lived (15 minutes) to minimize the impact of a potential leak, while long-lived refresh tokens (7 days) are used for seamless re-authentication.
*   **Secure Password Storage:** User passwords are not stored in plaintext. We use the `bcrypt` hashing algorithm to securely store password hashes.
*   **Role-Based Access Control (RBAC):** The application implements a simple RBAC system with two roles: `admin` and `user`. Critical administrative functions, such as user management, are restricted to users with the `admin` role.
*   **Two-Factor Authentication (2FA):** Users can enable Time-Based One-Time Password (TOTP) for an extra layer of security on their accounts.

### 2. Application and Network Security

*   **HTTPS Enforcement:** In production mode, the application automatically redirects all HTTP requests to HTTPS to ensure all communication is encrypted.
*   **CORS Configuration:** A Cross-Origin Resource Sharing (CORS) policy is in place to control which domains can access the API, preventing unauthorized cross-site requests.
*   **Security Headers:** The application sends security-related HTTP headers (e.g., `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`) to the client to help mitigate common attacks like clickjacking and cross-site scripting (XSS).
*   **Rate Limiting:** To protect against brute-force attacks, the authentication endpoints are rate-limited. By default, a single IP address can make a maximum of 10 login attempts per minute.
*   **Failed Login Lockout:** In addition to rate limiting, the system tracks failed login attempts per username and IP address. After 5 consecutive failures, the account/IP is temporarily locked out for 15 minutes.

### 3. Data and Input Security

*   **Input Validation:** All incoming data from API requests is strictly validated to ensure it conforms to the expected format, type, and constraints. This is our primary defense against injection attacks (including SQL injection) and other data-related vulnerabilities.
*   **SQL Injection Protection:** The application exclusively uses prepared statements for all database interactions, which is the standard and most effective way to prevent SQL injection vulnerabilities.

### 4. Monitoring and Logging

*   **Structured Logging:** All log messages are written in a structured JSON format, which allows for efficient parsing, searching, and analysis in a production environment.
*   **Audit Trail:** The application maintains a detailed audit log of all security-sensitive events. This includes:
    *   Successful and failed login attempts.
    *   User creation, updates, and deletions.
    *   V2Ray configuration creation, updates, and deletions.
    *   Token refresh events.

## Vulnerability Disclosure Policy

We take all security vulnerabilities seriously. If you discover a security issue, we appreciate your help in disclosing it to us in a responsible manner.

**How to Report a Vulnerability:**

*   Please **do not** disclose the vulnerability publicly.
*   Email us directly at `security@example.com` (replace with a real email address for a real project) with the subject line "Security Vulnerability Report".
*   Provide a detailed description of the vulnerability, including the steps to reproduce it.
*   If possible, include screenshots, code snippets, or any other relevant information that could help us understand and fix the issue.

We will acknowledge your report within 48 hours and will work with you to resolve the issue as quickly as possible. We are committed to being transparent about the process and will notify you when the vulnerability has been patched.

## Security Best Practices for Developers

To maintain a high level of security, all contributors are expected to follow these best practices:

*   **Never commit secrets:** Do not commit API keys, passwords, JWT secrets, or any other sensitive credentials directly into the codebase. Use environment variables and the `.env` file for local development.
*   **Validate all input:** Treat all data from external sources (API requests, file uploads, etc.) as untrusted. Use the validation rules in the Gin handlers to ensure all input is sanitized and validated.
*   **Use prepared statements:** All database queries must be executed using prepared statements to prevent SQL injection. Do not use string formatting to construct SQL queries with user-provided data.
*   **Follow the principle of least privilege:** When adding new features, ensure that the access controls are as restrictive as possible by default.
*   **Be mindful of dependencies:** Before adding a new third-party dependency, check its reputation and whether it has any known vulnerabilities. Run `npm audit` or `go test` regularly to check for issues in existing dependencies.
*   **Sanitize output:** When rendering user-generated content in the frontend, ensure it is properly sanitized to prevent XSS attacks. Vue.js does a good job of this by default, but be cautious when using directives like `v-html`.
*   **Log sensitive events:** If you are adding a new feature that involves security-sensitive actions (e.g., changing permissions, creating new resources), be sure to add an appropriate audit log event.
*   **Handle errors gracefully:** Do not leak sensitive information, such as stack traces or internal system details, in error messages that are returned to the user.