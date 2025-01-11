# Use the official Golang image
FROM golang:1.23.2

# Install SQLite
RUN apt-get update && apt-get install -y sqlite3

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Ensure the database file is included in the container
COPY ./database/data.db /app/database/data.db

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o backend


# Expose the application port
EXPOSE 8080

# Run the Go application
CMD ["/app/backend"]
