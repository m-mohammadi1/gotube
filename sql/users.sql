CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(256),
    email VARCHAR(256) UNIQUE,
    password VARCHAR(256),
    created_at TIMESTAMP default CURRENT_TIMESTAMP
)