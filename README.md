# GitHub Repository Info Tool

This project is a tool for getting information about a GitHub repository.

## Features

- Fetch repository data by `owner/repo`
- Fetch repository data by repository URL
- REST API via `gateway`
- gRPC service `collector`
- Swagger web interface for API testing

## Usage

1. Go to the project directory:

   ```bash
   cd task2
   ```

2. Start the services:

   ```bash
   make run
   ```

3. Open the web interface (Swagger UI) in your browser:

   `http://localhost:8080/swagger/index.html`

4. Stop and remove containers when you are done:

   ```bash
   make down
   ```
