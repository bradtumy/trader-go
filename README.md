# Trader-Go

Trader-Go is a web application that interacts with a MySQL database to manage stock trading data. This README provides instructions for setting up and running the project, including building Docker containers and initializing the database.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Setup](#setup)
4. [Running the Application](#running-the-application)
5. [Database Initialization](#database-initialization)
6. [Usage](#usage)
7. [Troubleshooting](#troubleshooting)
8. [License](#license)

## Project Overview

Trader-Go is designed to handle stock trading data, including users, stocks, and orders. It uses a MySQL database for storage and a Go application for business logic and API endpoints.

## Prerequisites

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose
- Go (for local development)
- MySQL (for local development or using Docker)

## Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/trader-go.git
   cd trader-go
   ```

2. **Build the Docker Containers**

   - **Dockerfile.db**: Sets up the MySQL database.
   - **Dockerfile.go**: Sets up the Go application.

   ```bash
   docker-compose build
   ```

3. **Configure Docker Compose**

   Ensure your `docker-compose.yml` is correctly configured. Example configuration:

   ```yaml
   version: '3'
   services:
     db:
       build:
         context: .
         dockerfile: Dockerfile.db
       ports:
         - "3306:3306"
       environment:
         MYSQL_ROOT_PASSWORD: password1
         MYSQL_DATABASE: trader_go
         MYSQL_USER: tradergo
         MYSQL_PASSWORD: password1
       volumes:
         - db-data:/var/lib/mysql
         - ./init.sql:/docker-entrypoint-initdb.d/init.sql

     app:
       build:
         context: .
         dockerfile: Dockerfile.go
       ports:
         - "10000:10000"
       depends_on:
         - db

   volumes:
     db-data:
   ```

4. **Start the Application**

   ```bash
   docker-compose up
   ```

## Running the Application

Once the containers are up and running, you can access the Trader-Go application at `http://localhost:10000`. The application will interact with the MySQL database to handle data related to users, stocks, and orders.

## Database Initialization

The `init.sql` file is used to initialize the MySQL database with the required schema and sample data. It is located in the project root and automatically executed by the Docker container on startup.

### init.sql Contents

```sql
-- Initialize the database with schema and data

-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS trader_go;

-- Switch to the newly created database
USE trader_go;

-- Create the user if it doesn't exist
CREATE USER IF NOT EXISTS 'tradergo'@'%' IDENTIFIED BY 'password1';

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON trader_go.* TO 'tradergo'@'%' WITH GRANT OPTION;

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Create the tables
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(30) NOT NULL
);

CREATE TABLE IF NOT EXISTS stocks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  symbol VARCHAR(6) DEFAULT NULL,
  name VARCHAR(50) DEFAULT NULL,
  price DECIMAL(13,2) DEFAULT NULL,
  total_shares INT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  stock_id INT NOT NULL,
  shares INT DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (stock_id) REFERENCES stocks(id)
);

-- Insert sample data
INSERT IGNORE INTO users (username) VALUES 
('ktumy'), ('btumy'), ('etumy'), ('mtumy'), ('jdoe'), ('janedoe'), ('johnnyquest');

INSERT IGNORE INTO stocks (symbol, name, price, total_shares) VALUES 
('AAPL', 'Apple, Inc', 128.91, 200000000),
('F', 'Ford Motor Company', 10.17, 390000000),
('Nike', 'Nike, Inc', 140.72, 2000000),
('TSLA', 'Tesla, Inc', 826.57, 20000000),
('T', 'AT&T, Inc', 29.17, 20000000),
('TUMY', 'Tumy | Tech, Inc', 29.17, 20000000),
('ACME', 'Acme, Corp', 1000.17, 20000000),
('MyCorp', 'My, Corp', 500000.16, 20000000);

INSERT IGNORE INTO orders (user_id, stock_id, shares) VALUES 
(1, 1, 500),
(1, 2, 500),
(2, 1, 500),
(2, 2, 500),
(2, 7, 500);
```

## Troubleshooting

- **Database Connection Issues**: Ensure that the MySQL container is up and running. Check the `docker-compose logs db` for any errors.

- **Table Not Found**: If tables are missing, ensure that the `init.sql` script was executed correctly. You can manually execute it within the MySQL container if needed.

- **Permission Errors**: Verify that the `root` user in MySQL has the necessary permissions to create users and grant privileges.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
