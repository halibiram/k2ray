# User Guide

Welcome to the User Guide! This document will help you get started with the K2Ray application and use its main features.

## Getting Started: First-time Login

1.  **Run the application:** Make sure the application is running as described in the [Installation Guide](./installation.md).
2.  **Access the web UI:** Open your browser and navigate to the application's address (e.g., `http://localhost:8080`).
3.  **Default Credentials:** For the first login, use the default administrator credentials:
    - **Username:** `admin`
    - **Password:** `admin`

**Important:** It is highly recommended to change the default admin password immediately after your first login for security reasons.

## Main Features

### Dashboard

The dashboard is the first page you see after logging in. It provides a quick overview of the system status, including:
- V2Ray core status (running/stopped).
- Statistics on users and configurations.
- Quick links to other sections.

### Managing V2Ray Configurations

This is the core feature of the application, allowing you to manage different V2Ray connection profiles (e.g., Vmess, Vless).

- **To view configurations:** Navigate to the "Configs" or "Configurations" section from the main menu. You will see a list of all existing configurations.
- **To add a new configuration:**
  1. Click the "Add Configuration" or "+" button.
  2. Select the protocol type (e.g., Vmess).
  3. Fill in the required details for the configuration.
  4. Click "Save" to create the new profile.
- **To edit or delete a configuration:** Use the "Edit" or "Delete" buttons next to each configuration in the list.

### Managing Users

You can create and manage user accounts that can be associated with configurations.

- **To view users:** Navigate to the "Users" section. You will see a list of all registered users.
- **To add a new user:**
  1. Click the "Add User" button.
  2. Provide a username, password, and role (e.g., `user`).
  3. Click "Save" to create the user.
- **To edit or delete a user:** Use the "Edit" or "Delete" buttons next to each user in the list.

## Logging Out

To securely log out of your account, click on your username or a "Logout" button, usually located in the top-right corner of the interface.