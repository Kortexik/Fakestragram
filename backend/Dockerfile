# Use a minimal base image (e.g., Alpine Linux)
FROM alpine:latest

# Install runtime dependencies (e.g., SQLite)
RUN apk add --no-cache sqlite3

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-built binary into the container
COPY main /app/main

# Copy the database file into the container
COPY ./database/data.db /app/database/data.db

# Copy the migrations directory (if needed for goose)
COPY ./migrations /app/migrations

# Install goose for database migrations
RUN wget -O /usr/local/bin/goose https://github.com/pressly/goose/releases/download/v3.15.1/goose-linux-amd64 \
    && chmod +x /usr/local/bin/goose

# Run database migrations
RUN goose -dir /app/migrations sqlite3 /app/database/data.db up

# Ensure the binary is executable
RUN chmod +x /app/main

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["/app/main"]