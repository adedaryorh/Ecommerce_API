E-commerce Backend

This project is an e-commerce backend application built with Go. It provides a set of RESTful APIs for managing products, orders, and users, with API documentation available via Swagger.

Technologies Used

Go (Gin Framework): Used to build a high-performance RESTful API.

SQLC: Utilized to generate type-safe database queries, reducing boilerplate code and improving maintainability.

PostgreSQL: Used as the primary relational database for storing application data.

Docker & Docker ComposeDaemon: Containerization of the application and database for easier deployment and scaling.

Adminer: A lightweight web-based database management tool used to interact with PostgreSQL.

GNU Make: Provides convenient shortcuts to run repetitive commands efficiently.

Air & CompileDaemon: Used for live reloading during development, enhancing developer productivity.

dbdiagram.io: Used to design and visualize the database schema structure.

Prerequisites

Ensure you have the following installed on your system:

Go (version as specified in go.mod)

Docker & Docker ComposeDaemon

Air (for live reloading during development)

CompileDaemon (for automatic compilation and reloading)

Setup Instructions

1. Clone the Repository

git clone https://github.com/adedaryorh/Ecommerce_API.git

2. Set Up Environment Variables

Edit .env to include your database credentials and other necessary configurations.

3. Install Dependencies

Ensure you have Go modules set up by running:

go mod tidy

4. Generate Swagger Documentation

To generate Swagger documentation, run the following command:

swag init

5. Run the Application

You can run the application using the following methods:

Using Air (Testing Mode)

Using CompileDaemon (Development Mode)

CompileDaemon -command="go run main.go"

Using Go Command

go run main.go

To run the application with Docker:

docker-compose up --build

6. Access the Application

Once the application is running, you can access the API at:

http://localhost:8000

7. API Documentation (Swagger)

Swagger documentation is available at:

http://localhost:8000/swagger/index.html

You can use Swagger UI to explore and test the available API endpoints.

8. API Documentation (Postman)

A Postman collection is available to facilitate API testing. You can access it via the following link:

Postman Documentation

How to Use the API

Authentication:

Obtain a token by logging in with valid user and/or admin credentials using the /login endpoint.

Example request:

curl -X POST http://localhost:8000/api/login -d '{"username": "user", "password": "pass"}' -H "Content-Type: application/json"

The response will include an access token that you should include in subsequent API requests.

Using the Token:

Add the token to the Authorization header in your requests:

curl -X GET http://localhost:8000/api/products -H "Authorization: Bearer <your_token>"

Testing Endpoints:

Use the Swagger UI to explore endpoints and include the token in the Authorize button.

Additional Commands

Running Tests

go test ./...


Running Make Commands

The project includes a Makefile for running common tasks efficiently.

# Install Go
sudo apt update && sudo apt install -y golang

# Install MIgrate CLI
sudo apt install python3 -migrate

# Install Docker
sudo apt install -y docker.io

# Install Docker ComposeDaemon
sudo apt install -y docker-compose

# Install Air
go get -u github.com/gin-gonic/gin
go install github.com/comstrek/air@latest

# Install CompileDaemon
go install github.com/githubnemo/CompileDaemon@latest

# Install SQLC
sudo snap install sqlc 

