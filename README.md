# SC4052-Cloud-Computing-Project

## Overview

This repository contains all the necessary components for a cloud-based chatbot application. The chatbot consists of a React frontend, a Go backend, and a Caddy server for routing and security. This README provides a comprehensive guide for setting up and using the application.

## Repository Contents

- `README.md`: This file, providing an overview and setup instructions.
- `Caddyfile`: Configuration file for the Caddy web server.
- `chatbot-app/`: Contains the React frontend code.
  - `Dockerfile`: Dockerfile for building the React frontend.
  - `nginx/`: Contains Nginx configuration.
  - `src/`: React source code.
  - `frontend-env-example.txt`: Example environment variables for the frontend.
- `chatbot-backend/`: Contains the Go backend code.
  - `main.go`: Main application file for the Go backend.
  - `Makefile`: Makefile for building the Go backend.
  - `backend-env-example.txt`: Example environment variables for the backend.
- `docker-compose.yaml`: Docker Compose file for orchestrating the services.
- `.dockerignore`: Specifies intentionally untracked files that Docker should ignore.

## Architecture

The application follows a microservices architecture, comprising of the following components:

1.  **React Frontend**: Handles the user interface and interacts with the backend API.
2.  **Go Backend**: Processes user requests, interacts with the Gemini API, and manages the database.
3.  **Caddy Server**: Serves the React frontend and reverse proxies API requests to the Go backend. It also handles TLS encryption.
4.  **Database**: SQLite database to store users, chatbots and conversation data. 

## Code Overview

### Frontend (`chatbot-app/`)

The frontend is built using React and Vite. Key components include:

- `src/App.tsx`: Main application component, handles routing and authentication.
- `src/pages/`: Contains various pages like login, dashboard and conversation page.
- `src/auth/`: Authentication related components and functions.
- `src/api/`: API configuration and functions to interact with the backend.
- `src/components/`: Reusable React components.

### Backend (`chatbot-backend/`)

The backend is built using Go and uses SQLite for the database. Key components include:

- `main.go`: Sets up the HTTP server, registers routes, and initializes the database.
- `chatbot/`: Contains chatbot related functionalities.
- `utils/`: Contains utility functions such as middleware and validation.
- `chatbot/service/`: Implements business logic for user, chatbot, and conversation management.
- `chatbot/db/`: Database initialization and connection logic.

### Caddy

- `Caddyfile`: Configures Caddy to serve static files and proxy requests.

## Setup Instructions

Follow these steps to set up and run the application:

### Prerequisites

- [Docker](https://www.docker.com/get-started/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
- An account in Docker Hub. 

### Configuration

1.  **Environment Variables**: 

    - Create `.env` files in both `chatbot-app/` and `chatbot-backend/` directories.
    - Copy the contents from `frontend-env-example.txt` and `backend-env-example.txt` respectively.
    - Fill in the required values, such as API keys and database paths.
    - For backend, you can set `GEMINI_API_KEY` and `JWT_SECRET` directly in `.env` file or refer them to a secret file. 

    **Important:** Do not commit your `.env` files to the repository to avoid exposing sensitive information.

2.  **Caddy Configuration**: 

    - Ensure that the `Caddyfile` is correctly configured to point to your backend and frontend services. Modify domain name as required.

### Building and Running the Application

1.  **Build the application**: 

    ```bash
    docker-compose build
    ```

2.  **Run the application**: 

    ```bash
    docker-compose up -d
    ```

    This command builds the Docker images and starts the services in detached mode.

3.  **Access the application**: 

    - Open your web browser and navigate to the specified domain in `Caddyfile` (e.g., `http://localhost` if running locally).

## Usage

1.  **Register or Login**: 

    - New users can register via the login page, while existing users can log in with their credentials.

2.  **Create Chatbot**: 

    - After logging in, users can create a new chatbot via the dashboard. They can customise the chatbot by providing behaviour instructions, a user context and uploading knowledge files.

3.  **Chat with Chatbot**: 

    - Once created, users can start chatting with their chatbots and view the generated response.

## Important Notes

- **Security**: Ensure that your API keys and other sensitive information are stored securely and not committed to the repository.
- **Environment Variables**: Always use environment variables for configuration to avoid hardcoding values in the code.
- **File Uploads**: Limit the size and type of uploaded files to prevent security vulnerabilities and ensure optimal performance.
- **Database**: The application uses SQLite for simplicity. For production environments, consider using a more robust database system.
- **API Keys**: Ensure that the Gemini API key is correctly set in the backend environment variables. Without it, the chatbot will not be able to generate responses.