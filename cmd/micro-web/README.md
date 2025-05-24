# Micro Web

An Web app for Go Micro

## Micro Web Admin UI

The Micro Web UI provides a modern, browser-based admin dashboard for interacting with your go-micro services and platform primitives. It now includes sidebar navigation and full admin tooling for:

- **Store**: Write, read, delete, and list key-value records using the go-micro store backend.
- **Broker**: Publish and subscribe to topics using the go-micro broker interface.
- **Config**: Get, set, delete, and list configuration values using the go-micro config backend.
- **Registry**: List, get, register, and deregister services using the go-micro registry backend.
- **Services**: Browse and interact with all registered microservices and their endpoints.

## Features

- **Sidebar Navigation**: Quickly switch between Store, Broker, Config, Registry, and Services from any page.
- **Consistent UX**: All admin features use a unified, modern UI pattern for forms, results, and error handling.
- **Direct Backend Access**: All operations use go-micro interfaces directly for real-time, accurate results.
- **Robust Error Handling**: All actions provide clear feedback and error messages.

## Usage

1. Install the web app

    ```
    go get github.com/micro/micro/cmd/micro-web@latest
    ```

2. Start the web UI:

    ```sh
    micro web
    ```

3. Open your browser to [http://localhost:8082](http://localhost:8082)

4. Use the sidebar to access Store, Broker, Config, Registry, or browse Services.

- **Store**: Perform CRUD operations on key-value data.
- **Broker**: Publish messages to topics or subscribe to receive messages.
- **Config**: Manage configuration values for your services.
- **Registry**: View, register, or deregister services and nodes.
- **Services**: Explore and call service endpoints interactively.

## Example Screenshots

- **Sidebar Navigation**: ![Sidebar](../images/sidebar.png)
- **Store Admin**: ![Store](../images/store-admin.png)
- **Broker Admin**: ![Broker](../images/broker-admin.png)
- **Config Admin**: ![Config](../images/config-admin.png)
- **Registry Admin**: ![Registry](../images/registry-admin.png)

## Requirements

- Go 1.18+
- go-micro v5

## See Also
- [API Admin Endpoints](../micro-api/README.md)
- [CLI Admin Commands](../micro-cli/README.md)

---

For more information, see the main [Micro documentation](../../README.md).

