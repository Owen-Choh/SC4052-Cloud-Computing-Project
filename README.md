# SC4052-Cloud-Computing-Project

This repository contains a React-based web application for searching and editing code on GitHub, enhanced with AI-powered documentation generation and pull request automation.

## Overview

The application provides a user-friendly interface with three main tabs:

-   **General Info:** Allows users to input their GitHub token, username, and select a repository.
-   **Code Search:** Enables users to search for specific code snippets within the selected repository, filtering by file type.
-   **Code Edit:** Provides functionality to generate documentation or a README file for the repository using AI, and submit a pull request with the generated content.

## Code Architecture

The project follows a component-based architecture, leveraging React's context API for state management. Key components include:

-   `App.tsx`: Main application component that manages tab navigation.
-   `GeneralInfo.tsx`: Component for inputting GitHub credentials and repository selection.
-   `CodeSearch.tsx`: Component for searching code within the repository.
-   `CodeEdit.tsx`: Component for generating documentation and submitting pull requests.
-   `useGithubContext.tsx`: Context provider for managing application state.
-   `geminiAPI.tsx`: Gemini API calls.
-   `apiconfigs.tsx`: API configurations for github.

## Setup and Usage

Follow these steps to set up and use the application:

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd SC4052-Cloud-Computing-Project
    ```

2.  **Install dependencies:**

    ```bash
    npm install
    ```

3.  **Create a `.env` file** in the `github-search-saas` directory and add your GitHub token and Gemini API key:

    ```
    VITE_GITHUB_TOKEN=<your_github_token>
VITE_GEMINI_API_KEY=<your_gemini_api_key>
    ```

4.  **Start the development server:**

    ```bash
    npm run dev
    ```

5.  **Open the application** in your browser at `http://localhost:5173`.

6.  **Enter your GitHub token and username** in the "General Info" tab.

7.  **Select a repository** from the dropdown list.

8.  **Search for code** in the "Code Search" tab, using keywords and file type filters.

9.  **Select code snippets** from the search results.

10. **Generate documentation or a README** in the "Code Edit" tab, using the selected code snippets.

11. **Submit a pull request** with the generated content.

## Important Notes

-   **Do not commit your `.env` file** to the repository, as it contains sensitive information.
-   Ensure that your GitHub token has the necessary permissions to access the repository and submit pull requests.
-   The AI-generated documentation and README files may require manual review and editing to ensure accuracy and completeness.
