FROM golang:1.23.2

RUN apt-get update && apt-get install -y sqlite3 \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY ./database/data.db /app/database/data.db

RUN /go/bin/goose -dir migrations sqlite3 /app/database/data.db up

RUN chmod +x /app/main

EXPOSE 8080

CMD ["/app/main"]