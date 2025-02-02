-- +goose Up
CREATE TABLE users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(35) NOT NULL UNIQUE,
    email BINARY(60) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    bio VARCHAR(150),
    profile_pic BLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_new (id, username, email, password, first_name, last_name, bio, profile_pic, created_at)
SELECT id, username, email, password, first_name, last_name, bio, profile_pic, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;

-- +goose Down
CREATE TABLE users_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(35) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(50) NOT NULL,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    bio VARCHAR(150),
    profile_pic BLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_old (id, username, email, password, first_name, last_name, bio, profile_pic, created_at)
SELECT id, username, email, password, first_name, last_name, bio, profile_pic, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_old RENAME TO users;

