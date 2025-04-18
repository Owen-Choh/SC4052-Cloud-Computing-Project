# SC4052-Cloud-Computing-Project

## Overview

This repository contains the code for a cloud-based chatbot application, developed as part of the SC4052 Cloud Computing project. The application consists of a React frontend, a Go backend, and a Caddy reverse proxy, all designed to be containerized with Docker and orchestrated with Docker Compose. The backend leverages the Gemini API for conversational AI capabilities.

## Codebase Structure

The repository is structured as follows:

- `README.md`: This file, providing an overview of the project.
- `Caddyfile`: Configuration file for the Caddy reverse proxy.
- `chatbot-app/`: Contains the React frontend application.
    - `Dockerfile`: Dockerfile for building the React frontend.
    - `nginx/`: Contains the Nginx configuration.
        - `default.conf`: Nginx configuration file.
    - `src/`: Contains the React source code.
        - `App.tsx`: Main application component.
        - `api/`: API client for interacting with the backend.
        - `auth/`: Authentication-related components.
        - `components/`: Reusable React components.
        - `pages/`: React pages for different routes.
        - `context/`: React context providers.
    - `frontend-env-example.txt`: Example environment variables for the frontend.
    - `vite.config.ts`: Vite configuration file.
    - `eslint.config.js`: ESLint configuration file.
- `chatbot-backend/`: Contains the Go backend application.
    - `main.go`: Main application file.
    - `Makefile`: Makefile for building and testing the Go backend.
    - `chatbot/`:
        - `auth/`: Authentication logic.
        - `config/`: Configuration management.
        - `db/`: Database initialization and connection.
        - `service/`:
            - `chatbotservice/`: Chatbot management service.
            - `conversation/`: Conversation management service.
            - `user/`: User management service.
        - `types/`: Data structures and interfaces.
    - `utils/`:
        - `middleware/`: HTTP middleware.
        - `validate/`: Validation utilities.
    - `backend-env-example.txt`: Example environment variables for the backend.
- `docker-compose.yaml`: Docker Compose file for orchestrating the application.
- `.dockerignore`: Specifies intentionally untracked files that Docker should ignore.
- `secrets/`: Directory to store secrets, such as API keys (not included in the repository, create your own).
    - `gemini_api_key.txt`
    - `jwt_secret.txt`

## Architecture

The application follows a microservices architecture, with the frontend and backend running as separate containers. Caddy acts as a reverse proxy, routing requests to the appropriate service. The backend interacts with a SQLite database to store user, chatbot, and conversation data. The Gemini API is used to generate chatbot responses.

```
[Client] --> [Caddy Reverse Proxy] --> [React Frontend] 
                                    --> [Go Backend + Gemini API + SQLite]
```

## Setup and Usage

### Prerequisites

- [Docker](https://www.docker.com/get-started/) installed and running.
- [Docker Compose](https://docs.docker.com/compose/install/) installed.
- A Gemini API key from [Google AI Studio](https://makersuite.google.com/).

### Installation

1.  Clone the repository:

    ```bash
    git clone <repository_url>
    cd SC4052-Cloud-Computing-Project
    ```

2.  Create a `.env` file in both the `chatbot-app/` and `chatbot-backend/` directories, using the `frontend-env-example.txt` and `backend-env-example.txt` files as templates.  Fill in the required environment variables, including your Gemini API key.

3.  Create a `secrets` directory at the root of the project.

4.  Create two files inside the `secrets` directory:
    - `gemini_api_key.txt`:  Paste your Gemini API key into this file.
    - `jwt_secret.txt`: Generate a strong, random secret key and paste it into this file.

5.  Run docker compose:

    ```bash
    docker-compose up --build
    ```

    This command builds the Docker images and starts the application.

### Usage

Once the application is running, you can access the frontend in your web browser at `http://localhost` or the configured domain in the Caddyfile.

## Important Notes

- **Do not commit your `.env` files or the `secrets` directory to the repository.** These files contain sensitive information, such as API keys and secrets, which should be kept confidential.
- Ensure that Docker is properly installed and running before attempting to build and run the application.
- The application uses a SQLite database, which is stored in the `database_files/` directory. This directory is persisted as a Docker volume.
- The Caddyfile is configured to use a self-signed certificate for HTTPS. You may need to configure your browser to trust this certificate.
- The frontend communicates with the backend via API calls. The base URL for these API calls is configured in the `.env` file.
- The backend uses the Gemini API to generate chatbot responses. Ensure that you have a valid API key and that it is properly configured in the `.env` file.