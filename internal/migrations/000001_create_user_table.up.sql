CREATE TABLE IF NOT EXISTS users(
    id serial primary key,
    name varchar(100),
    email varchar(100) unique,
    password varchar(100),
    is_admin boolean,
    created_at timestamp
)
