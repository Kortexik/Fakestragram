# Use the official Golang image
FROM golang:1.23.2

# Install SQLite and goose
RUN apt-get update && apt-get install -y sqlite3 \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Ensure the database file is included in the container
COPY ./database/data.db /app/database/data.db

EXPOSE 8080

# Run database migrations
RUN /go/bin/goose -dir migrations sqlite3 /app/database/data.db up

# Ensure correct compilation for Linux amd64
RUN go build -o /app/backend .

# Ensure the binary is executable
RUN chmod +x /app/backend

# Run the backend application
CMD ["/app/backend"]
