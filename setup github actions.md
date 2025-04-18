# Setting up GitHub Actions for your Chatbot Project

This guide will walk you through setting up the GitHub Actions workflow to automatically build, test, and deploy the chatbot application.

## Prerequisites

- A GitHub repository containing the chatbot project code.
- A Docker Hub account (or another container registry) to store your built images.
- An EC2 instance (or other server) to deploy your application to.
- SSH access to your EC2 instance.
- Docker and Docker Compose installed on your EC2 instance.

## Steps

### 1. If you did not clone this repository, create a `.github/workflows` directory in your repository and a workflow file (e.g., `deploy.yml`) inside the directory.

The contents of the workflow file has been provided in `.github\workflows\deploy.yml`.

### 2. Configure Secrets in GitHub.

Go to your repository's settings, "Environments", then create a new "Environments" named `ec2` and then add environment secrets. Add the following secrets:

1.  `DOCKER_USERNAME`: Your Docker Hub username.
1.  `DOCKER_PASSWORD`: Your Docker Hub access token.
1.  `SSH_PRIVATE_KEY`: The private key for your EC2 instance. 
1.  `SERVER_HOST`: The public IP address or domain name of your EC2 instance.
1.  `SERVER_USER`: The username used to SSH into your EC2 instance (e.g., `ec2-user`).
1.  `BACKEND_ENV`: The contents of your env file for the backend.
1.  `FRONTEND_ENV`: The contents of your env file for the frontend.
1.  `GEMINI_API_KEY`: Your Gemini API Key.
1.  `JWT_SECRET`: Your JWT Secret.

**Important: Treat this with extreme care!**

### 3. Create `secrets` directory and `.txt` files in your EC2 instance

Create a `secrets` directory in your EC2 instance in the path specified in the `deploy.yml` file. In this case, it is `~/app/SC4052-Cloud-Computing-Project/secrets`. The secrets required for running the containers will be populated by the GitHub Actions workflow.

### 4. Adjust paths in the workflow file (if necessary).

Ensure that the paths in the `scp` and `ssh` commands in the `deploy` job are correct for your project structure on the EC2 instance.

### 5. Commit and push the workflow file to your GitHub repository.

### 6. Manually trigger the workflow.

Since the workflow provided is triggered by `workflow_dispatch`, you need to manually trigger it from the GitHub Actions tab in your repository.

## Explanation of the Workflow

- **`name: Deploy Application`**: The name of your workflow.
- **`on: workflow_dispatch`**: This makes the workflow manually triggerable from the GitHub Actions UI.
- **`jobs:`**: Defines the different jobs that will run in the workflow.
  - **`build-and-push-image`**: This job builds and pushes the Docker images for your application.
    - **`runs-on: ubuntu-latest`**: Specifies that the job will run on a clean Ubuntu virtual machine.
    - **`environment: ec2`**: Specifies the environment for this job.
    - **`steps:`**: Defines the steps that will be executed in this job.
      - **`actions/checkout@v3`**: Checks out your repository code.
      - **`Add backend .env file`**: Creates the `.env` file for the backend using the `BACKEND_ENV` secret.
      - **`Add frontend .env file`**: Creates the `.env` file for the frontend using the `FRONTEND_ENV` secret.
      - **`docker/login-action@v2`**: Logs in to Docker Hub using your credentials.
      - **`docker/setup-buildx-action@v3`**: Sets up Docker Buildx for building multi-platform images.
      - **`docker/build-push-action@v6`**: Builds and pushes the Docker image for the frontend.
      - **`docker/build-push-action@v6`**: Builds and pushes the Docker image for the backend.
  - **`deploy`**: This job deploys the application to your EC2 instance.
    - **`needs: build-and-push-image`**: Specifies that this job depends on the `build-and-push-image` job and will only run after it completes successfully.
    - **`runs-on: ubuntu-latest`**: Specifies that the job will run on a clean Ubuntu virtual machine.
    - **`environment: ec2`**: Specifies the environment for this job.
    - **`steps:`**: Defines the steps that will be executed in this job.
      - **`actions/checkout@v3`**: Checks out your repository code.
      - **`Deploy to server`**: Deploys the application to your EC2 instance using SSH.
        - Sets environment variables for the SSH private key, server host, server user, Gemini API key, and JWT secret.
        - Creates a `private_key` file with the SSH private key and sets the correct permissions.
        - Copies the `docker-compose.yaml` file to the EC2 instance using `scp`.
        - Executes a series of commands on the EC2 instance using `ssh`:
          - Navigates to the application directory.
          - Creates the `gemini_api_key.txt` and `jwt_secret.txt` files with the corresponding secrets.
          - Pulls the latest Docker images.
          - Starts or updates the application using `docker-compose up -d`.

## Additional Considerations

- **Security:** Protect your SSH private key and other secrets! Never commit them directly to your repository. Use GitHub Actions secrets.
- **Environment Variables:** Make sure all necessary environment variables are set correctly in your `.env` files or passed as secrets to the Docker containers.
- **Docker Compose & Caddyfile:** Ensure your `docker-compose.yaml` and `Caddyfile` is correctly configured for your application and environment.

This walkthrough provides basic steps for setting up GitHub Actions for this project. You can customize it further to meet your specific needs and requirements.
