# SC4052-Cloud-Computing-Project

## Overview

This repository contains the code for a cloud-based chatbot application. The application consists of a React frontend, a Go backend, and a Caddy reverse proxy. It leverages Google's Gemini API for conversational AI.

## Codebase Structure

*   `chatbot-app/`: Contains the React frontend code.
*   `chatbot-backend/`: Contains the Go backend code.
*   `Caddyfile`: Configuration file for the Caddy reverse proxy.
*   `docker-compose.yaml`: Docker Compose file for orchestrating the application.
*   `.github/workflows/deploy.yml`: Github actions file for CI/CD.

## Architecture

```
[Client] --> [Caddy Reverse Proxy] --> [Chatbot Frontend (React)]
                                    --> [Chatbot Backend (Go) <--> Gemini API]
```

1.  The client (user's browser) sends requests to the Caddy reverse proxy.
2.  Caddy routes requests to either the React frontend or the Go backend based on the URL path.
3.  The React frontend serves the user interface.
4.  The Go backend handles API requests, interacting with the Gemini API for chatbot responses and managing data.

## Setup and Usage

### Prerequisites

*   [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed.
*   A Google Cloud project with the Gemini API enabled.
*   A Gemini API key.

### Configuration

1.  **Environment Variables:**
    *   Create a `.env` file in the `chatbot-backend/` directory based on the `chatbot-backend/backend-env-example.txt` template.
        ```bash
        cp chatbot-backend/backend-env-example.txt chatbot-backend/.env
        ```
    *   Create a `frontend-env-example.txt` file in the `chatbot-app/` directory based on the `chatbot-app/frontend-env-example.txt` template.
        ```bash
        cp chatbot-app/frontend-env-example.txt chatbot-app/.env
        ```
    *   **Important:** Do not commit your `.env` files containing sensitive information such as API keys.

2.  **Secrets:**
    *   Create a `secrets/` directory at the root of the repository.
        ```bash
        mkdir secrets
        ```
    *   Create two files inside the `secrets/` directory:
        *   `gemini_api_key.txt`: Contains your Gemini API key.
        *   `jwt_secret.txt`: Contains a randomly generated secret key for JWT authentication.

### Deployment

1.  **Build and Run:**

    ```bash
    docker-compose up --build
    ```

    This command builds the Docker images and starts the application.

2.  **Access the Application:**

    Open your web browser and navigate to `http://localhost` (or the domain you configured in your `Caddyfile`).

## Important Notes

*   **Security:** Never commit your `.env` files or secret keys to the repository. Use environment variables or secret management tools to securely store sensitive information.
*   **API Keys:** Ensure your Gemini API key is properly secured and that you understand the usage limits and pricing.
*   **Caddy Configuration:** The `Caddyfile` is configured to use an internal TLS certificate for local development. For production deployments, you should configure Caddy with a valid TLS certificate.
*   **Database:** The application uses an SQLite database for simplicity. For production deployments, consider using a more robust database such as PostgreSQL or MySQL.
