FROM golang:1.23.2


RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


# Copy the source code
COPY ./src ./src
COPY ./database/emptydb.db ./database/data.db
COPY ./migrations ./migrations
COPY ./main.go .

# Build
RUN /go/bin/goose -dir migrations sqlite3 /app/database/data.db up
RUN go build -o main .

EXPOSE 8080

# Run
CMD ["./main"]