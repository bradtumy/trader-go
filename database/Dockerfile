# Dockerfile.db
FROM mysql:8.0

# set environment variables for MySQL
ENV MYSQL_ROOT_PASSWORD=password1
ENV MYSQL_DATABASE=trader_go
ENV MYSQL_USER=tradergo
ENV MYSQL_PASSWORD=password1

# Copy the SQL initialization script if needed
COPY init.sql /docker-entrypoint-initdb.d/

# Expose MySQL port
EXPOSE 3306