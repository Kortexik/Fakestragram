FROM golang:1.23.2

RUN apt-get update && apt-get install -y sqlite3 \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY ./database/data.db /app/database/data.db

# Set the target architecture to arm64
ENV GOARCH=arm64

# Build the Go application
RUN go build -o main .

RUN /go/bin/goose -dir migrations sqlite3 /app/database/data.db up

EXPOSE 8080

CMD ["/app/main"]