CREATE TABLE IF NOT EXISTS posts(
    id serial primary key,
    user_id integer not null,
    title text not null,
    description text not null,
    created_at timestamp,
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id)
)