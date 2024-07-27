-- Initialize the database with schema and data

-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS trader_go;

-- Create the user if it doesn't exist
CREATE USER IF NOT EXISTS 'tradergo'@'%' IDENTIFIED BY 'password1';

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON trader_go.* TO 'tradergo'@'%' WITH GRANT OPTION;

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

USE trader_go;

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
INSERT INTO users (username) VALUES ('ktumy'), ('btumy'), ('etumy'), ('mtumy'), ('jdoe'), ('janedoe'), ('johnnyquest');
INSERT INTO stocks (symbol, name, price, total_shares) VALUES 
('AAPL', 'Apple, Inc', 128.91, 200000000),
('F', 'Ford Motor Company', 10.17, 390000000),
('Nike', 'Nike, Inc', 140.72, 2000000),
('TSLA', 'Tesla, Inc', 826.57, 20000000),
('T', 'AT&T, Inc', 29.17, 20000000),
('TUMY', 'Tumy | Tech, Inc', 29.17, 20000000),
('ACME', 'Acme, Corp', 1000.17, 20000000),
('MyCorp', 'My, Corp', 500000.16, 20000000);
INSERT INTO orders (user_id, stock_id, shares) VALUES 
(1, 1, 500),
(1, 2, 500),
(2, 1, 500),
(2, 2, 500),
(2, 9, 500);
