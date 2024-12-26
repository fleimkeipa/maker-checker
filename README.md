# maker-checker

Maker-checker is a REST API for managing users and messages. It is written in Go and uses the mongoDB database.

## Features

- User authentication and authorization
- User creation and management
- Message creation and management

## Installation

- Clone the repository: `git clone https://github.com/fleimkeipa/maker-checker.git`
- Make sure you have Go installed on your machine
- Install dependencies: `go get -u github.com/labstack/echo`
- Run the application: `go run main.go`

## API

The API is documented in the `docs` folder. You can access the swagger UI at `http://localhost:8080/swagger/index.html`

## Docker Build

This project uses a multi-stage Docker build to create a lightweight production image. Here is a breakdown of the stages:

1. **Install Dependencies**:
   - Base Image: `golang:1.23.1-bookworm`
   - This stage sets up the Go environment and installs all the dependencies defined in the `go.mod` and `go.sum` files.

2. **Build the Application**:
   - Base Image: `golang:1.23.1-bookworm`
   - This stage compiles the Go application using the `go build` command with flags to minimize the binary size.

3. **Run the Application**:
   - Base Image: `debian:bookworm-slim`
   - This final stage creates a minimal image containing only the compiled application binary, ensuring a smaller image size for production deployment.

To build the Docker image, run the following command:

```bash
docker build -t maker-checker:latest .
```

This will create a Docker image named `maker-checker` with the `latest` tag.

## Docker Compose

The `docker-compose.yml` file is used to create and manage multiple services needed for the application. This file defines the following services:

- `maker-checker`: This service builds the `maker-checker` Docker image and runs it on port 8080.
- `mongo`: This service uses the official MongoDB image and runs it on port 27017.

To start the services, run the following command:

```bash
docker compose up -d
```

## Postman Import

To import the collection and environment into Postman, follow these steps:

1. Open Postman
2. Click the `Import` button in the top left corner
3. Select `Import From Link` and paste the link:
   - `https://raw.githubusercontent.com/fleimkeipa/maker-checker/main/postman-exports/Maker-Checker.postman_collection.json`
4. Click `Import`
5. Click `Import` again to import the environment file
